package service

import "github.com/ydataai/azure-quota-provider/pkg/common"

// RESTServiceConfiguration defines required configuration for rest service
type RESTServiceConfiguration struct {
	location    string
	machineType string
}

// LoadEnvVars parses the required configuration variables. Throws an error if the validations aren't met
func (c *RESTServiceConfiguration) LoadEnvVars() error {
	location, err := common.VariableFromEnvironment("LOCATION")
	if err != nil {
		return err
	}
	c.location = location

	machineType, err := common.VariableFromEnvironment("MACHINE_TYPE")
	if err != nil {
		return err
	}
	c.machineType = machineType

	return nil
}
