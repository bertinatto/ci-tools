// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a VPC endpoint service to which service consumers (Amazon Web Services
// accounts, users, and IAM roles) can connect.
//
// Before you create an endpoint service, you must create one of the following for
// your service:
//
//   - A [Network Load Balancer]. Service consumers connect to your service using an interface endpoint.
//
//   - A [Gateway Load Balancer]. Service consumers connect to your service using a Gateway Load Balancer
//     endpoint.
//
// If you set the private DNS name, you must prove that you own the private DNS
// domain name.
//
// For more information, see the [Amazon Web Services PrivateLink Guide].
//
// [Gateway Load Balancer]: https://docs.aws.amazon.com/elasticloadbalancing/latest/gateway/
// [Network Load Balancer]: https://docs.aws.amazon.com/elasticloadbalancing/latest/network/
// [Amazon Web Services PrivateLink Guide]: https://docs.aws.amazon.com/vpc/latest/privatelink/
func (c *Client) CreateVpcEndpointServiceConfiguration(ctx context.Context, params *CreateVpcEndpointServiceConfigurationInput, optFns ...func(*Options)) (*CreateVpcEndpointServiceConfigurationOutput, error) {
	if params == nil {
		params = &CreateVpcEndpointServiceConfigurationInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateVpcEndpointServiceConfiguration", params, optFns, c.addOperationCreateVpcEndpointServiceConfigurationMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateVpcEndpointServiceConfigurationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateVpcEndpointServiceConfigurationInput struct {

	// Indicates whether requests from service consumers to create an endpoint to your
	// service must be accepted manually.
	AcceptanceRequired *bool

	// Unique, case-sensitive identifier that you provide to ensure the idempotency of
	// the request. For more information, see [How to ensure idempotency].
	//
	// [How to ensure idempotency]: https://docs.aws.amazon.com/ec2/latest/devguide/ec2-api-idempotency.html
	ClientToken *string

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// The Amazon Resource Names (ARNs) of the Gateway Load Balancers.
	GatewayLoadBalancerArns []string

	// The Amazon Resource Names (ARNs) of the Network Load Balancers.
	NetworkLoadBalancerArns []string

	// (Interface endpoint configuration) The private DNS name to assign to the VPC
	// endpoint service.
	PrivateDnsName *string

	// The supported IP address types. The possible values are ipv4 and ipv6 .
	SupportedIpAddressTypes []string

	// The Regions from which service consumers can access the service.
	SupportedRegions []string

	// The tags to associate with the service.
	TagSpecifications []types.TagSpecification

	noSmithyDocumentSerde
}

type CreateVpcEndpointServiceConfigurationOutput struct {

	// Unique, case-sensitive identifier that you provide to ensure the idempotency of
	// the request.
	ClientToken *string

	// Information about the service configuration.
	ServiceConfiguration *types.ServiceConfiguration

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateVpcEndpointServiceConfigurationMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpCreateVpcEndpointServiceConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpCreateVpcEndpointServiceConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CreateVpcEndpointServiceConfiguration"); err != nil {
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
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateVpcEndpointServiceConfiguration(options.Region), middleware.Before); err != nil {
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

func newServiceMetadataMiddleware_opCreateVpcEndpointServiceConfiguration(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CreateVpcEndpointServiceConfiguration",
	}
}
