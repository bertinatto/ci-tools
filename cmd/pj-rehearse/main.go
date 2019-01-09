package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	// TODO: Solve this properly
	"github.com/getlantern/deepcopy"
	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"

	pjapi "k8s.io/test-infra/prow/apis/prowjobs/v1"
	pjclientset "k8s.io/test-infra/prow/client/clientset/versioned"
	pj "k8s.io/test-infra/prow/client/clientset/versioned/typed/prowjobs/v1"
	prowconfig "k8s.io/test-infra/prow/config"
	"k8s.io/test-infra/prow/pjutil"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/ci-operator-prowgen/pkg/diffs"
)

func getPrNumber(jobSpec *pjapi.ProwJobSpec) int {
	return jobSpec.Refs.Pulls[0].Number
}

func getJobSpec() (*pjapi.ProwJobSpec, error) {
	specEnv := []byte(os.Getenv("JOB_SPEC"))
	if len(specEnv) == 0 {
		return nil, fmt.Errorf("JOB_SPEC not set or set to an empty string")
	}
	spec := pjapi.ProwJobSpec{}
	if err := json.Unmarshal(specEnv, &spec); err != nil {
		return nil, err
	}

	if len(spec.Refs.Pulls) > 1 {
		return nil, fmt.Errorf("Cannot rehearse in the context of a batch job")
	}

	return &spec, nil
}

func loadClusterConfig() (*rest.Config, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err == nil {
		return clusterConfig, nil
	}

	credentials, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, fmt.Errorf("could not load credentials from config: %v", err)
	}

	clusterConfig, err = clientcmd.NewDefaultClientConfig(*credentials, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load client configuration: %v", err)
	}
	return clusterConfig, nil
}

// Rehearsed Prow jobs may depend on ConfigMaps with content also modified by
// the tested PR. All ci-operator-based jobs use the `ci-operator-configs`
// ConfigMap that contains ci-operator configuration files. Rehearsed jobs
// need to have the PR-version of these files available. The following code
// takes care of creating a short-lived, rehearsal ConfigMap. The keys needed
// to be present are extracted from the rehearsal jobs and the rehearsal jobs
// are modified to use this ConfigMap instead of the "production" one.

var ciOperatorConfigsCMName = "ci-operator-configs"

type rehearsalCIOperatorConfigs struct {
	cmclient  corev1.ConfigMapInterface
	prNumber  int
	configDir string

	logger logrus.FieldLogger
	dry    bool

	configMapName string
	neededConfigs map[string]string
}

func newRehearsalCIOperatorConfigs(cmclient corev1.ConfigMapInterface, prNumber int, configDir string, logger logrus.FieldLogger, dry bool) *rehearsalCIOperatorConfigs {
	name := fmt.Sprintf("rehearsal-ci-operator-configs-%d", prNumber)
	return &rehearsalCIOperatorConfigs{
		cmclient:      cmclient,
		prNumber:      prNumber,
		configDir:     configDir,
		logger:        logger.WithField("ciop-configs-cm", name),
		dry:           dry,
		configMapName: name,
		neededConfigs: map[string]string{},
	}
}

// If a job uses the `ci-operator-config` ConfigMap, save which key does it use
// from it and replace that ConfigMap reference with a reference to the
// temporary, rehearsal ConfigMap containing the necessary keys with content
// matching the version from tested PR
func (c *rehearsalCIOperatorConfigs) FixupJob(job *prowconfig.Presubmit, repo string) {
	for _, container := range job.Spec.Containers {
		for _, env := range container.Env {
			if env.ValueFrom == nil {
				continue
			}
			if env.ValueFrom.ConfigMapKeyRef == nil {
				continue
			}
			if env.ValueFrom.ConfigMapKeyRef.Name == ciOperatorConfigsCMName {
				filename := env.ValueFrom.ConfigMapKeyRef.Key
				env.ValueFrom.ConfigMapKeyRef.Name = c.configMapName
				c.neededConfigs[filename] = filepath.Join(repo, filename)

				logFields := logrus.Fields{"ci-operator-config": filename, "rehearsal-job": job.Name}
				c.logger.WithFields(logFields).Info("Rehearsal job uses ci-operator config ConfigMap")
			}
		}
	}
}

// Create a rehearsal ConfigMap with ci-operator config files needed by the
// rehearsal jobs.
func (c *rehearsalCIOperatorConfigs) Create() error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: c.configMapName},
		Data:       map[string]string{},
	}
	c.logger.Debug("Preparing rehearsal ConfigMap for ci-operator configs")

	for key, path := range c.neededConfigs {
		fullPath := filepath.Join(c.configDir, path)
		content, err := ioutil.ReadFile(fullPath)
		c.logger.WithField("ciop-config", key).Info("Loading ci-operator config to rehearsal ConfigMap")
		if err != nil {
			return fmt.Errorf("failed to read ci-operator config file from %s: %v", fullPath, err)
		}

		cm.Data[key] = string(content)
	}

	if c.dry {
		cmAsYAML, err := yaml.Marshal(cm)
		if err != nil {
			return fmt.Errorf("Failed to marshal ConfigMap to YAML: %v", err)
		}
		fmt.Printf("%s\n", cmAsYAML)
		return nil
	}
	c.logger.Info("Creating rehearsal ConfigMap for ci-operator configs")
	_, err := c.cmclient.Create(cm)
	return err
}

func makeRehearsalPresubmit(source *prowconfig.Presubmit, repo string, prNumber int) (*prowconfig.Presubmit, error) {
	var rehearsal prowconfig.Presubmit
	deepcopy.Copy(&rehearsal, source)

	rehearsal.Name = fmt.Sprintf("rehearse-%d-%s", prNumber, source.Name)
	rehearsal.Context = fmt.Sprintf("ci/rehearse/%s/%s", repo, strings.TrimPrefix(source.Context, "ci/prow/"))

	if len(source.Spec.Containers) != 1 {
		return nil, fmt.Errorf("Cannot rehearse jobs with more than 1 container in Spec")
	}
	container := source.Spec.Containers[0]

	if len(container.Command) != 1 || container.Command[0] != "ci-operator" {
		return nil, fmt.Errorf("Cannot rehearse jobs that have Command different from simple 'ci-operator'")
	}

	for _, arg := range container.Args {
		if strings.HasPrefix(arg, "--git-ref") || strings.HasPrefix(arg, "-git-ref") {
			return nil, fmt.Errorf("Cannot rehearse jobs that call ci-operator with '--git-ref' arg")
		}
	}

	if len(source.Branches) != 1 {
		return nil, fmt.Errorf("Cannot rehearse jobs that run over multiple branches")
	}
	branch := strings.TrimPrefix(strings.TrimSuffix(source.Branches[0], "$"), "^")

	gitrefArg := fmt.Sprintf("--git-ref=%s@%s", repo, branch)
	rehearsal.Spec.Containers[0].Args = append(source.Spec.Containers[0].Args, gitrefArg)

	return &rehearsal, nil
}

func submitRehearsal(job *prowconfig.Presubmit, jobSpec *pjapi.ProwJobSpec, logger logrus.FieldLogger, pjclient pj.ProwJobInterface, dry bool) (*pjapi.ProwJob, error) {
	labels := make(map[string]string)
	for k, v := range job.Labels {
		labels[k] = v
	}

	pj := pjutil.NewProwJob(pjutil.PresubmitSpec(*job, *(jobSpec.Refs)), labels)
	logger.WithFields(pjutil.ProwJobFields(&pj)).Info("Submitting a new prowjob.")

	if dry {
		jobAsYAML, err := yaml.Marshal(pj)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal job to YAML: %v", err)
		}
		fmt.Printf("%s\n", jobAsYAML)
		return &pj, nil
	}

	return pjclient.Create(&pj)
}

func execute(toBeRehearsed map[string][]prowconfig.Presubmit, jobSpec *pjapi.ProwJobSpec, logger logrus.FieldLogger, rehearsalConfigs *rehearsalCIOperatorConfigs, pjclient pj.ProwJobInterface, dry bool) error {
	rehearsals := []*prowconfig.Presubmit{}

	for repo, jobs := range toBeRehearsed {
		for _, job := range jobs {
			jobLogger := logger.WithFields(logrus.Fields{"target-repo": repo, "target-job": job.Name})
			rehearsal, err := makeRehearsalPresubmit(&job, repo, getPrNumber(jobSpec))
			if err != nil {
				jobLogger.WithError(err).Warn("Failed to make a rehearsal presubmit")
			} else {
				jobLogger.WithField("rehearsal-job", rehearsal.Name).Info("Created a rehearsal job to be submitted")
				rehearsalConfigs.FixupJob(rehearsal, repo)
				rehearsals = append(rehearsals, rehearsal)
			}
		}
	}

	if len(rehearsals) > 0 {
		if err := rehearsalConfigs.Create(); err != nil {
			return fmt.Errorf("failed to prepare rehearsal ci-operator config ConfigMap: %v", err)
		}
		for _, job := range rehearsals {
			created, err := submitRehearsal(job, jobSpec, logger, pjclient, dry)
			if err != nil {
				logger.WithError(err).Warn("Failed to execute a rehearsal presubmit presubmit")
			} else {
				logger.WithFields(pjutil.ProwJobFields(created)).Info("Submitted rehearsal prowjob")
			}
		}
	} else {
		logger.Warn("No job rehearsals")
	}

	return nil
}

type options struct {
	dryRun bool

	configPath    string
	jobConfigPath string

	ciopConfigsPath string
}

func gatherOptions() options {
	o := options{}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	fs.BoolVar(&o.dryRun, "dry-run", true, "Whether to actually submit rehearsal jobs to Prow")

	fs.StringVar(&o.configPath, "config-path", "/etc/config/config.yaml", "Path to Prow config.yaml")
	fs.StringVar(&o.jobConfigPath, "job-config-path", "", "Path to prow job configs.")

	fs.StringVar(&o.ciopConfigsPath, "ci-operator-configs", "", "Path to a directory containing ci-operator configs")

	fs.Parse(os.Args[1:])
	return o
}

func validateOptions(o options) error {
	if len(o.jobConfigPath) == 0 {
		return fmt.Errorf("empty --job-config-path")
	}
	return nil
}

func main() {
	o := gatherOptions()
	err := validateOptions(o)
	if err != nil {
		logrus.WithError(err).Fatal("invalid options")
	}

	jobSpec, err := getJobSpec()
	if err != nil {
		logrus.WithError(err).Fatal("could not read JOB_SPEC")
	}

	prFields := logrus.Fields{"org": jobSpec.Refs.Org, "repo": jobSpec.Refs.Repo, "PR": getPrNumber(jobSpec)}
	logger := logrus.WithFields(prFields)
	logger.Info("Rehearsing Prow jobs for a configuration PR")

	prowConfig, err := prowconfig.Load(o.configPath, o.jobConfigPath)
	if err != nil {
		logger.WithError(err).Fatal("Failed to load Prow config")
	}
	prowjobNamespace := prowConfig.ProwJobNamespace

	clusterConfig, err := loadClusterConfig()
	if err != nil {
		logger.WithError(err).Fatal("could not load cluster clusterConfig")
	}

	pjcset, err := pjclientset.NewForConfig(clusterConfig)
	if err != nil {
		logger.WithError(err).Fatal("could not create a ProwJob clientset")
	}
	pjclient := pjcset.ProwV1().ProwJobs(prowjobNamespace)

	cmcset, err := corev1.NewForConfig(clusterConfig)
	if err != nil {
		logger.WithError(err).Fatal("could not create a Core clientset")
	}
	cmclient := cmcset.ConfigMaps(prowjobNamespace)

	rehearsalConfigs := newRehearsalCIOperatorConfigs(cmclient, getPrNumber(jobSpec), o.ciopConfigsPath, logger, o.dryRun)

	changedPresubmits, err := diffs.GetChangedPresubmits(o.configPath, o.jobConfigPath)
	if err != nil {
		logger.WithError(err).Fatal("Failed to determine which jobs should be rehearsed")
	}

	if err := execute(changedPresubmits, jobSpec, logger, rehearsalConfigs, pjclient, o.dryRun); err != nil {
		logger.WithError(err).Fatal("Failed to execute rehearsal jobs")
	}
}

// TODO: Validate o.ciopConfigsPath
// TODO: Migrate GetChangedPresubmits to accept full config
// TODO: Use shared version of loadClusterConfig
// TODO: Extract job handling / rehearsal CM stuff to a package
