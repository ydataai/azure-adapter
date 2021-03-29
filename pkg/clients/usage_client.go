package clients

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"

	"github.com/sirupsen/logrus"
)

// UsageClientInterface defines a interface for usage client
type UsageClientInterface interface {
	ComputeUsage(context.Context, string, string) (compute.Usage, error)
}

// UsageClient defines a struct with required dependencies for usage client
type UsageClient struct {
	logger *logrus.Logger
	client compute.UsageClient
}

// NewUsageClient initializes usage client
func NewUsageClient(
	logger *logrus.Logger,
	client compute.UsageClient,
) UsageClient {
	return UsageClient{
		logger: logger,
		client: client,
	}
}

// ComputeUsage fetches compute usage list and filter a machine type
func (uc UsageClient) ComputeUsage(ctx context.Context, location string, machineType string) (compute.Usage, error) {
	uc.logger.Info("Starting to fetch available GPU in Azure")

	result, err := uc.client.List(ctx, location)
	if err != nil {
		uc.logger.Infof("while fetching usage result. Error: %v", err)
		return compute.Usage{}, err
	}

	var gpuUsage compute.Usage
	for _, usage := range result.Values() {
		if *usage.Name.Value == machineType {
			gpuUsage = usage
			break
		}
	}

	return gpuUsage, nil
}
