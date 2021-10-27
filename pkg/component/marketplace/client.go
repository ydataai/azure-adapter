package marketplace

import "github.com/Azure/go-autorest/autorest"

const baseURI = "https://marketplaceapi.microsoft.com/api"

// BaseClient is the base client for Compute.
type BaseClient struct {
	autorest.Client
	BaseURI        string
	SubscriptionID string
}

// New creates an instance of the BaseClient client.
func New(subscriptionID string) BaseClient {
	return NewWithBaseURI(baseURI, subscriptionID)
}

// NewWithBaseURI creates an instance of the BaseClient client using a custom endpoint.  Use this when interacting with
// an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return BaseClient{
		Client:         autorest.NewClientWithOptions(autorest.ClientOptions{}),
		BaseURI:        baseURI,
		SubscriptionID: subscriptionID,
	}
}
