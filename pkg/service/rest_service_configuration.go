package service

import (
	"github.com/kelseyhightower/envconfig"
)

// RESTServiceConfiguration defines required configuration for rest service
type RESTServiceConfiguration struct {
	Location    string `envconfig:"LOCATION" required:"true"`
	MachineType string `envconfig:"MACHINE_TYPE" required:"true"`
}

// LoadFromEnvVars parses the required configuration variables
func (c *RESTServiceConfiguration) LoadFromEnvVars() error {
	return envconfig.Process("", c)
}
