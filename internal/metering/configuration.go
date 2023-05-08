// Package metering provides objects to interact with metering API
package metering

import "github.com/kelseyhightower/envconfig"

// TimeLayout ISO time layout
const TimeLayout = "2006-01-02T15:04:05.000Z"

// Configuration represents the configuration for metering client.
type Configuration struct {
	ResourceUri string `envconfig:"MANAGED_APP_RESOURCE_URI" required:"true"`
	PlanId      string `envconfig:"MANAGED_APP_PLAN_ID" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the metering client.
func (c *Configuration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
