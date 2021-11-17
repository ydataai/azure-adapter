package marketplace

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/services/cloud"
)

const APIVersion = "2018-08-31"

// MarketplaceClient defines a struct with required dependencies for marketplace client
type MarketplaceClient struct {
	logger logging.Logger
	config AzureMarketplaceConfiguration
	client BaseClient
}

// NewMarketplaceClient initializes marketplace client
func NewMarketplaceClient(config AzureMarketplaceConfiguration, authorizer autorest.Authorizer, logger logging.Logger) cloud.MeteringService {
	client := New()
	client.Authorizer = authorizer
	return &MarketplaceClient{
		config: config,
		logger: logger,
		client: client,
	}
}

// CreateUsageEvent sends usage event to marketplace api for metering purpose.
func (c MarketplaceClient) CreateUsageEvent(ctx context.Context, event cloud.UsageEventReq) (cloud.UsageEventRes, error) {
	azevent := UsageEventReq{
		ResourceId: c.config.resourceUri,
		Plan:       c.config.planId,
		Dimension:  event.DimensionID,
		StartTime:  event.StartAt,
		Quantity:   event.Quantity,
	}

	req, err := c.CreateUsageEventPreparer(ctx, azevent)
	if err != nil {
		return cloud.UsageEventRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEvent", "CreateUsageEvent", nil, "Failure preparing request")
	}
	res, err := c.CreateSender(req)
	if err != nil {
		return cloud.UsageEventRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEvent", "CreateUsageEvent", res, "Failure responding to request")
	}
	azres, err := c.CreateUsageEventResponder(res)
	if err != nil {
		return cloud.UsageEventRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEvent", "CreateUsageEvent", res, "Failure creating response")
	}

	return cloud.UsageEventRes{
		UsageEventID: azres.UsageEventId,
		DimensionID:  azres.ResourceId,
		Status:       azres.Status,
	}, nil
}

// CreateUsageEventBatch sends usage batch events to marketplace api for metering purpose.
func (c MarketplaceClient) CreateUsageEventBatch(ctx context.Context, batch cloud.UsageEventBatchReq) (cloud.UsageEventBatchRes, error) {
	events := []UsageEventReq{}
	resourceDimension := map[string]string{}
	for _, request := range batch.Request {
		resourceDimension[c.config.resourceUri] = request.DimensionID
		event := UsageEventReq{
			ResourceId: c.config.resourceUri,
			Plan:       c.config.planId,
			Dimension:  request.DimensionID,
			StartTime:  request.StartAt,
			Quantity:   request.Quantity,
		}
		events = append(events, event)
	}

	req, err := c.CreateUsageEventBatchPreparer(ctx, UsageEventBatchReq{Request: events})
	if err != nil {
		return cloud.UsageEventBatchRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEventBatch", "CreateUsageEventBatch", nil, "Failure preparing request")
	}
	res, err := c.CreateSender(req)
	if err != nil {
		return cloud.UsageEventBatchRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEventBatch", "CreateUsageEventBatch", res, "Failure responding to request")
	}
	azres, err := c.CreateUsageEventBatchResponder(res)
	if err != nil {
		return cloud.UsageEventBatchRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEventBatch", "CreateUsageEventBatch", res, "Failure creating response")
	}

	results := []cloud.UsageEventRes{}
	for _, result := range azres.Result {
		if len(result.Error.Details) > 0 {
			c.logger.Errorf("Failed to process event batch %v.", result.Error.Details)
		}
		event := cloud.UsageEventRes{
			UsageEventID: result.UsageEventId,
			DimensionID:  result.Dimension,
			Status:       result.Status,
		}
		results = append(results, event)
	}
	return cloud.UsageEventBatchRes{Result: results}, nil
}

// CreatePreparer prepares the Create request.
func (c MarketplaceClient) CreateUsageEventPreparer(ctx context.Context, req UsageEventReq) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.client.BaseURI),
		autorest.WithPath("/usageEvent"),
		autorest.WithJSON(req),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateResponder handles the response to the UsageEventRes request. The method always
// closes the http.Response Body.
func (c MarketplaceClient) CreateUsageEventResponder(resp *http.Response) (UsageEventRes, error) {
	result := UsageEventRes{}
	err := autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	if err != nil {
		return result, err
	}
	return result, nil
}

// CreateUsageEventBatchPreparer prepares the batch request.
func (c MarketplaceClient) CreateUsageEventBatchPreparer(ctx context.Context, req UsageEventBatchReq) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.client.BaseURI),
		autorest.WithPath("/batchUsageEvent"),
		autorest.WithJSON(req),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateUsageEventBatchResponder handles the response to the UsageEventBAtchRes request. The method always
// closes the http.Response Body.
func (c MarketplaceClient) CreateUsageEventBatchResponder(resp *http.Response) (UsageEventBatchRes, error) {
	result := UsageEventBatchRes{}
	err := autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	if err != nil {
		return result, err
	}
	return result, nil
}

// CreateSender sends the request. The method will close the
// http.Response Body if it receives an error.
func (c MarketplaceClient) CreateSender(req *http.Request) (*http.Response, error) {
	return c.client.Send(req, azure.DoRetryWithRegistration(c.client.Client))
}
