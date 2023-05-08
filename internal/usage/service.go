// Package usage offers objects and methods to help using usage APIs
package usage

import (
	"context"

	"github.com/ydataai/go-core/pkg/common/logging"
)

const (
	vCPUToGPUFactor int64 = 6
)

// RESTServiceInterface defines rest service interface
type RESTService interface {
	AvailableGPU(ctx context.Context) (GPU, error)
}

// RESTService defines a struct with required dependencies for rest service
type restService struct {
	logger        logging.Logger
	configuration RESTServiceConfiguration
	usageClient   Client
}

// NewRESTService initializes rest service
func NewRESTService(
	logger logging.Logger,
	configuration RESTServiceConfiguration,
	usageClient Client,
) RESTService {
	return restService{
		logger:        logger,
		configuration: configuration,
		usageClient:   usageClient,
	}
}

// AvailableGPU ..
func (rs restService) AvailableGPU(ctx context.Context) (GPU, error) {
	rs.logger.Infof("fetch available GPUs from quota API")

	usageResult, err := rs.usageClient.ComputeUsage(ctx, rs.configuration.Location, rs.configuration.MachineType)
	if err != nil {
		rs.logger.Errorf(
			"while fetching list of %s/%s with error %v", rs.configuration.Location, rs.configuration.MachineType, err)
		return GPU(0), err
	}

	rs.logger.Infof(
		"resources for %s/%s -> Limit: %d | Current: %d",
		rs.configuration.Location, rs.configuration.MachineType, *usageResult.Limit, *usageResult.CurrentValue)

	availableGPU := (*usageResult.Limit - int64(*usageResult.CurrentValue)) / vCPUToGPUFactor

	rs.logger.Infof("Number of GPUs available in %s: %d", rs.configuration.Location, availableGPU)

	return GPU(availableGPU), nil
}
