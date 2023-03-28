// Package usage offers objects and methods to help using usage APIs
package usage

import (
	"context"

	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

// Client defines an interface for usage client
type Client interface {
	ComputeUsage(context.Context, string, string) (compute.Usage, error)
}

type usageClient struct {
	client *compute.UsageClient
}

// NewClient initializes an usage client
func NewClient(client *compute.UsageClient) Client {
	return usageClient{
		client: client,
	}
}

// ComputeUsage fetches compute usage list and filter a machine type
func (c usageClient) ComputeUsage(ctx context.Context, location string, machineType string) (compute.Usage, error) {
	pager := c.client.NewListPager(location, nil)

	var gpuUsage compute.Usage

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			return compute.Usage{}, err
		}

		for _, usage := range nextResult.Value {
			if *usage.Name.Value == machineType {
				gpuUsage = *usage
				break
			}
		}
	}

	return gpuUsage, nil
}
