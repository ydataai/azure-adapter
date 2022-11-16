package marketplace

import "github.com/kelseyhightower/envconfig"

// TimeLayout ISO time layout
const TimeLayout = "2006-01-02T15:04:05.000Z"

// AzureMarketplaceConfiguration represents the configuration for marketplace client.
type AzureMarketplaceConfiguration struct {
	ResourceUri string `envconfig:"AZURE_MANAGED_APP_RESOURCE_URI" required:"true"`
	PlanId      string `envconfig:"AZURE_MANAGED_APP_PLAN_ID" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the marketplace client.
func (c *AzureMarketplaceConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
