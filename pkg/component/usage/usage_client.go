package usage

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
)

// UsageClientInterface defines a interface for usage client
type UsageClientInterface interface {
	ComputeUsage(context.Context, string, string) (compute.Usage, error)
}

// UsageClient defines a struct with required dependencies for usage client
type UsageClient struct {
	client compute.UsageClient
}

// NewUsageClient initializes usage client
func NewUsageClient(client compute.UsageClient) UsageClient {
	return UsageClient{
		client: client,
	}
}

// ComputeUsage fetches compute usage list and filter a machine type
func (uc UsageClient) ComputeUsage(ctx context.Context, location string, machineType string) (compute.Usage, error) {
	result, err := uc.client.List(ctx, location)
	if err != nil {
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
