// Package metering provides objects to interact with metering API
package metering

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/metering"
)

type apiPath string

const (
	usageEventAPIPath      apiPath = "usageEvent"
	batchUsageEventAPIPath apiPath = "batchUsageEvent"
)

const (
	apiVersion = "2018-08-31"
	baseURI    = "https://marketplaceapi.microsoft.com/api"
)

// Client defines a struct with required dependencies for metering client
type client struct {
	config Configuration
	logger logging.Logger
	pl     runtime.Pipeline
}

// NewClient initializes metering client
func NewClient(
	credential azcore.TokenCredential, config Configuration, logger logging.Logger,
) (metering.Client, error) {
	pl, err := armruntime.NewPipeline("marketplace", "v0.1.0", credential, runtime.PipelineOptions{}, nil)
	if err != nil {
		return client{}, err
	}

	return client{
		config: config,
		logger: logger,
		pl:     pl,
	}, nil
}

// CreateUsageEvent creates and sends a request to create an UsageEvent
// It returns an error if any or an UsageEventResponse from azure
func (c client) CreateUsageEvent(ctx context.Context, event metering.UsageEvent) (metering.UsageEventResponse, error) {
	c.logger.Infof("received create event with %+v", event)

	if event.Quantity <= 0 {
		c.logger.Infof("metric '%s' skipped (%s <-> %s) = %v",
			event.DimensionID,
			event.StartAt.Format(TimeLayout),
			time.Now().Format(TimeLayout),
			event.Quantity,
		)
		return metering.UsageEventResponse{}, nil
	}

	azevent := usageEvent{
		Dimension:          event.DimensionID,
		Quantity:           event.Quantity,
		EffectiveStartTime: event.StartAt,
		ResourceURI:        c.config.ResourceUri,
		PlanID:             c.config.PlanId,
	}

	c.logger.Infof("event transformed into %+v", azevent)

	req, err := createRequest(ctx, usageEventAPIPath, azevent)
	if err != nil {
		return metering.UsageEventResponse{}, err
	}

	resp, err := c.pl.Do(req)
	if err != nil {
		return metering.UsageEventResponse{}, err
	}

	c.logger.Infof("got response %+v", resp)
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("failed to decode body with error %v", err)
	}

	c.logger.Info("body: ", string(bytes))

	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return metering.UsageEventResponse{}, invalidStatusCodeError(resp)
	}

	eventResponse := usageEventResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &eventResponse); err != nil {
		return metering.UsageEventResponse{}, err
	}

	c.logger.Infof("unmarshelled into event %+v", eventResponse)

	return metering.UsageEventResponse{
		UsageEventID: eventResponse.UsageEventId,
		DimensionID:  eventResponse.Dimension,
		Status:       eventResponse.Status,
	}, nil
}

// CreateUsageEventBatch creates a batch of UsageEvent with azure APIs
// It returns an error if any or an UsageEventResponse from azure
func (c client) CreateUsageEventBatch(
	ctx context.Context, batch metering.UsageEventBatch,
) (metering.UsageEventBatchResponse, error) {
	events := []usageEvent{}

	for _, request := range batch.Events {
		if request.Quantity <= 0 {
			c.logger.Infof("metric '%s' skipped (%s <-> %s) = %v",
				request.DimensionID,
				request.StartAt.Format(TimeLayout),
				time.Now().Format(TimeLayout),
				request.Quantity,
			)
			continue
		}

		event := usageEvent{
			Dimension:          request.DimensionID,
			Quantity:           request.Quantity,
			EffectiveStartTime: request.StartAt,
			ResourceURI:        c.config.ResourceUri,
			PlanID:             c.config.PlanId,
		}
		events = append(events, event)
	}

	req, err := createRequest(ctx, batchUsageEventAPIPath, usageEventBatch{Events: events})
	if err != nil {
		return metering.UsageEventBatchResponse{}, err
	}

	resp, err := c.pl.Do(req)
	if err != nil {
		return metering.UsageEventBatchResponse{}, err
	}

	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return metering.UsageEventBatchResponse{}, invalidStatusCodeError(resp)
	}

	result := &usageEventBatchResponse{}
	if err := runtime.UnmarshalAsJSON(resp, result); err != nil {
		return metering.UsageEventBatchResponse{}, err
	}

	results := []metering.UsageEventResponse{}
	for _, result := range result.Result {
		if len(result.Error.Details) > 0 {
			c.logger.Errorf("Failed to process event batch %v.", result.Error.Details)
		}
		event := metering.UsageEventResponse{
			UsageEventID: result.UsageEventId,
			DimensionID:  result.Dimension,
			Status:       result.Status,
		}
		results = append(results, event)
	}
	return metering.UsageEventBatchResponse{Result: results}, nil
}

func createRequest(ctx context.Context, path apiPath, event interface{}) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(baseURI, string(path)))
	if err != nil {
		return nil, err
	}

	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}

	return req, runtime.MarshalAsJSON(req, event)
}

func invalidStatusCodeError(response *http.Response) error {
	return fmt.Errorf("request failed with error %s", response.Status)
}
