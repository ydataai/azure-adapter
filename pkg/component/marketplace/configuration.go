package marketplace

import "github.com/kelseyhightower/envconfig"

// AzureMarketplaceConfiguration represents the configuration for marketplace client.
type AzureMarketplaceConfiguration struct {
	resourceUri string `envconfig:"AZURE_MANAGED_APP_RESOURCE_URI" required:"true"`
	planId      string `envconfig:"AZURE_MANAGED_APP_PLAN_ID" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the marketplace client.
func (c *AzureMarketplaceConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
