package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ydataai/azure-quota-provider/pkg/common"
)

type configuration struct {
	logLevel       logrus.Level
	subscriptionID string
}

func (c *configuration) LoadEnvVars() error {
	logLevel, err := common.VariableFromEnvironment("LOG_LEVEL")
	if err != nil {
		return err
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	c.logLevel = level

	c.subscriptionID, err = common.VariableFromEnvironment("SUBSCRIPTION_ID")
	if err != nil {
		return err
	}

	return nil
}
