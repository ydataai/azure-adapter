package main

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration defines all env vars required for the application
type Configuration struct {
	SubscriptionID string `envconfig:"ARM_SUBSCRIPTION_ID" required:"true"`
}

// LoadEnvVars reads all env vars required for the server package
func (c *Configuration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
