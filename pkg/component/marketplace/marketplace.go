package marketplace

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/sirupsen/logrus"
)

const APIVersion = "2018-08-31"

// MarketplaceClient defines a struct with required dependencies for marketplace client
type MarketplaceClient struct {
	logger *logrus.Logger
	Client BaseClient
}

// NewMarketplaceClient initializes marketplace client
func NewMarketplaceClient(logger *logrus.Logger) MarketplaceClient {
	return MarketplaceClient{
		logger: logger,
		Client: New(),
	}
}

// CreateUsageEvent sends usage event to marketplace api for metering purpose.
func (c MarketplaceClient) CreateUsageEvent(ctx context.Context, event UsageEventReq) (UsageEventRes, error) {
	req, err := c.CreateUsageEventPreparer(ctx, event)
	if err != nil {
		return UsageEventRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEvent", "CreateUsageEvent", nil, "Failure preparing request")
	}
	res, err := c.CreateSender(req)
	if err != nil {
		return UsageEventRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEvent", "CreateUsageEvent", res, "Failure responding to request")
	}
	return c.CreateUsageEventResponder(res)
}

// CreateUsageEventBatch sends usage batch events to marketplace api for metering purpose.
func (c MarketplaceClient) CreateUsageEventBatch(ctx context.Context, events []UsageEventReq) (UsageEventBatchRes, error) {
	req, err := c.CreateUsageEventBatchPreparer(ctx, UsageEventBatchReq{Request: events})
	if err != nil {
		return UsageEventBatchRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEventBatch", "CreateUsageEventBatch", nil, "Failure preparing request")
	}
	res, err := c.CreateSender(req)
	if err != nil {
		return UsageEventBatchRes{}, autorest.NewErrorWithError(err, "marketplace.UsageEventBatch", "CreateUsageEventBatch", res, "Failure responding to request")
	}
	return c.CreateUsageEventBatchResponder(res)
}

// CreatePreparer prepares the Create request.
func (c MarketplaceClient) CreateUsageEventPreparer(ctx context.Context, req UsageEventReq) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.Client.BaseURI),
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
		autorest.WithBaseURL(c.Client.BaseURI),
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
	return c.Client.Send(req, azure.DoRetryWithRegistration(c.Client.Client))
}
