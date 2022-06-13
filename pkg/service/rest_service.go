package service

import (
	"context"

	"github.com/ydataai/azure-adapter/pkg/component/usage"

	"github.com/ydataai/go-core/pkg/common/logging"
)

const (
	vCPUToGPUFactor int64 = 6
)

// RESTServiceInterface defines rest service interface
type RESTServiceInterface interface {
	AvailableGPU(ctx context.Context) (usage.GPU, error)
}

// RESTService defines a struct with required dependencies for rest service
type RESTService struct {
	logger        logging.Logger
	configuration RESTServiceConfiguration
	usageClient   usage.UsageClientInterface
}

// NewRESTService initializes rest service
func NewRESTService(
	logger logging.Logger,
	configuration RESTServiceConfiguration,
	usageClient usage.UsageClientInterface,
) RESTService {
	return RESTService{
		logger:        logger,
		configuration: configuration,
		usageClient:   usageClient,
	}
}

// AvailableGPU ..
func (rs RESTService) AvailableGPU(ctx context.Context) (usage.GPU, error) {
	rs.logger.Infof("fetch available GPUs from quota API")

	usageResult, err := rs.usageClient.ComputeUsage(ctx, rs.configuration.Location, rs.configuration.MachineType)
	if err != nil {
		rs.logger.Errorf(
			"while fetching list of %s/%s with error %v", rs.configuration.Location, rs.configuration.MachineType, err)
		return usage.GPU(0), err
	}

	rs.logger.Infof(
		"got result for resources %s/%s: %v", rs.configuration.Location, rs.configuration.MachineType, usageResult)

	availableGPU := (*usageResult.Limit - int64(*usageResult.CurrentValue)) / vCPUToGPUFactor

	rs.logger.Infof("Number of GPUs available in %s: %d", rs.configuration.Location, availableGPU)

	return usage.GPU(availableGPU), nil
}
