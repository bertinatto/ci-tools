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

// Lists the resources from a resource scan. The results can be filtered by
// resource identifier, resource type prefix, tag key, and tag value. Only
// resources that match all specified filters are returned. The response indicates
// whether each returned resource is already managed by CloudFormation.
func (c *Client) ListResourceScanResources(ctx context.Context, params *ListResourceScanResourcesInput, optFns ...func(*Options)) (*ListResourceScanResourcesOutput, error) {
	if params == nil {
		params = &ListResourceScanResourcesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListResourceScanResources", params, optFns, c.addOperationListResourceScanResourcesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListResourceScanResourcesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListResourceScanResourcesInput struct {

	// The Amazon Resource Name (ARN) of the resource scan.
	//
	// This member is required.
	ResourceScanId *string

	//  If the number of available results exceeds this maximum, the response includes
	// a NextToken value that you can use for the NextToken parameter to get the next
	// set of results. By default the ListResourceScanResources API action will return
	// at most 100 results in each response. The maximum value is 100.
	MaxResults *int32

	// A string that identifies the next page of resource scan results.
	NextToken *string

	// If specified, the returned resources will have the specified resource
	// identifier (or one of them in the case where the resource has multiple
	// identifiers).
	ResourceIdentifier *string

	// If specified, the returned resources will be of any of the resource types with
	// the specified prefix.
	ResourceTypePrefix *string

	// If specified, the returned resources will have a matching tag key.
	TagKey *string

	// If specified, the returned resources will have a matching tag value.
	TagValue *string

	noSmithyDocumentSerde
}

type ListResourceScanResourcesOutput struct {

	// If the request doesn't return all the remaining results, NextToken is set to a
	// token. To retrieve the next set of results, call ListResourceScanResources
	// again and use that value for the NextToken parameter. If the request returns
	// all results, NextToken is set to an empty string.
	NextToken *string

	// List of up to MaxResults resources in the specified resource scan that match
	// all of the specified filters.
	Resources []types.ScannedResource

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListResourceScanResourcesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsquery_serializeOpListResourceScanResources{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpListResourceScanResources{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ListResourceScanResources"); err != nil {
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
	if err = addOpListResourceScanResourcesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListResourceScanResources(options.Region), middleware.Before); err != nil {
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

// ListResourceScanResourcesPaginatorOptions is the paginator options for
// ListResourceScanResources
type ListResourceScanResourcesPaginatorOptions struct {
	//  If the number of available results exceeds this maximum, the response includes
	// a NextToken value that you can use for the NextToken parameter to get the next
	// set of results. By default the ListResourceScanResources API action will return
	// at most 100 results in each response. The maximum value is 100.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListResourceScanResourcesPaginator is a paginator for ListResourceScanResources
type ListResourceScanResourcesPaginator struct {
	options   ListResourceScanResourcesPaginatorOptions
	client    ListResourceScanResourcesAPIClient
	params    *ListResourceScanResourcesInput
	nextToken *string
	firstPage bool
}

// NewListResourceScanResourcesPaginator returns a new
// ListResourceScanResourcesPaginator
func NewListResourceScanResourcesPaginator(client ListResourceScanResourcesAPIClient, params *ListResourceScanResourcesInput, optFns ...func(*ListResourceScanResourcesPaginatorOptions)) *ListResourceScanResourcesPaginator {
	if params == nil {
		params = &ListResourceScanResourcesInput{}
	}

	options := ListResourceScanResourcesPaginatorOptions{}
	if params.MaxResults != nil {
		options.Limit = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListResourceScanResourcesPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListResourceScanResourcesPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListResourceScanResources page.
func (p *ListResourceScanResourcesPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListResourceScanResourcesOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxResults = limit

	optFns = append([]func(*Options){
		addIsPaginatorUserAgent,
	}, optFns...)
	result, err := p.client.ListResourceScanResources(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// ListResourceScanResourcesAPIClient is a client that implements the
// ListResourceScanResources operation.
type ListResourceScanResourcesAPIClient interface {
	ListResourceScanResources(context.Context, *ListResourceScanResourcesInput, ...func(*Options)) (*ListResourceScanResourcesOutput, error)
}

var _ ListResourceScanResourcesAPIClient = (*Client)(nil)

func newServiceMetadataMiddleware_opListResourceScanResources(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ListResourceScanResources",
	}
}
