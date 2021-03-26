package service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/ydataai/azure-quota-provider/pkg/clients"
	"github.com/ydataai/azure-quota-provider/pkg/common"
)

const (
	vCPUToGPUFactor int64 = 6
)

// RESTServiceInterface defines rest service interface
type RESTServiceInterface interface {
	AvailableGPU(ctx context.Context) (common.GPU, error)
}

// RESTService defines a struct with required dependencies for rest service
type RESTService struct {
	logger        *logrus.Logger
	configuration RESTServiceConfiguration
	usageClient   clients.UsageClientInterface
}

// NewRESTService initializes rest service
func NewRESTService(
	logger *logrus.Logger,
	configuration RESTServiceConfiguration,
	usageClient clients.UsageClientInterface,
) RESTService {
	return RESTService{
		logger:        logger,
		configuration: configuration,
		usageClient:   usageClient,
	}
}

// AvailableGPU ..
func (rs RESTService) AvailableGPU(ctx context.Context) (common.GPU, error) {
	usageResult, err := rs.usageClient.ComputeUsage(ctx, rs.configuration.location, rs.configuration.machineType)
	if err != nil {
		rs.logger.Errorf("while fetching list of compute usage. Error %v", err)
		return common.GPU(0), err
	}

	availableGPU := (*usageResult.Limit - int64(*usageResult.CurrentValue)) / vCPUToGPUFactor

	return common.GPU(availableGPU), nil
}
