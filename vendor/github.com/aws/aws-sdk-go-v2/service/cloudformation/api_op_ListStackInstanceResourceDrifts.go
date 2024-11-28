// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudformation

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns drift information for resources in a stack instance.
//
// ListStackInstanceResourceDrifts returns drift information for the most recent
// drift detection operation. If an operation is in progress, it may only return
// partial results.
func (c *Client) ListStackInstanceResourceDrifts(ctx context.Context, params *ListStackInstanceResourceDriftsInput, optFns ...func(*Options)) (*ListStackInstanceResourceDriftsOutput, error) {
	if params == nil {
		params = &ListStackInstanceResourceDriftsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListStackInstanceResourceDrifts", params, optFns, c.addOperationListStackInstanceResourceDriftsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListStackInstanceResourceDriftsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListStackInstanceResourceDriftsInput struct {

	// The unique ID of the drift operation.
	//
	// This member is required.
	OperationId *string

	// The name of the Amazon Web Services account that you want to list resource
	// drifts for.
	//
	// This member is required.
	StackInstanceAccount *string

	// The name of the Region where you want to list resource drifts.
	//
	// This member is required.
	StackInstanceRegion *string

	// The name or unique ID of the stack set that you want to list drifted resources
	// for.
	//
	// This member is required.
	StackSetName *string

	// [Service-managed permissions] Specifies whether you are acting as an account
	// administrator in the organization's management account or as a delegated
	// administrator in a member account.
	//
	// By default, SELF is specified. Use SELF for stack sets with self-managed
	// permissions.
	//
	//   - If you are signed in to the management account, specify SELF .
	//
	//   - If you are signed in to a delegated administrator account, specify
	//   DELEGATED_ADMIN .
	//
	// Your Amazon Web Services account must be registered as a delegated
	//   administrator in the management account. For more information, see [Register a delegated administrator]in the
	//   CloudFormation User Guide.
	//
	// [Register a delegated administrator]: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacksets-orgs-delegated-admin.html
	CallAs types.CallAs

	// The maximum number of results to be returned with a single call. If the number
	// of available results exceeds this maximum, the response includes a NextToken
	// value that you can assign to the NextToken request parameter to get the next
	// set of results.
	MaxResults *int32

	// If the previous paginated request didn't return all of the remaining results,
	// the response object's NextToken parameter value is set to a token. To retrieve
	// the next set of results, call this action again and assign that token to the
	// request object's NextToken parameter. If there are no remaining results, the
	// previous response object's NextToken parameter is set to null .
	NextToken *string

	// The resource drift status of the stack instance.
	//
	//   - DELETED : The resource differs from its expected template configuration in
	//   that the resource has been deleted.
	//
	//   - MODIFIED : One or more resource properties differ from their expected
	//   template values.
	//
	//   - IN_SYNC : The resource's actual configuration matches its expected template
	//   configuration.
	//
	//   - NOT_CHECKED : CloudFormation doesn't currently return this value.
	StackInstanceResourceDriftStatuses []types.StackResourceDriftStatus

	noSmithyDocumentSerde
}

type ListStackInstanceResourceDriftsOutput struct {

	// If the previous paginated request didn't return all of the remaining results,
	// the response object's NextToken parameter value is set to a token. To retrieve
	// the next set of results, call this action again and assign that token to the
	// request object's NextToken parameter. If there are no remaining results, the
	// previous response object's NextToken parameter is set to null .
	NextToken *string

	// A list of StackInstanceResourceDriftsSummary structures that contain
	// information about the specified stack instances.
	Summaries []types.StackInstanceResourceDriftsSummary

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListStackInstanceResourceDriftsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpListStackInstanceResourceDrifts{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpListStackInstanceResourceDrifts{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ListStackInstanceResourceDrifts"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addOpListStackInstanceResourceDriftsValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListStackInstanceResourceDrifts(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opListStackInstanceResourceDrifts(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ListStackInstanceResourceDrifts",
	}
}
