// Package configuration provides objects to configure adapter objects
package configuration

import (
	"github.com/kelseyhightower/envconfig"
)

// Application defines all env vars required for the application
type Application struct {
	SubscriptionID string `envconfig:"AZURE_SUBSCRIPTION_ID" required:"true"`
}

// LoadFromEnvVars reads all env vars required for the server package
func (c *Application) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
