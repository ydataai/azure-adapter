package main

import (
	"github.com/ydataai/azure-adapter/pkg/common"

	"github.com/sirupsen/logrus"
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

	c.subscriptionID, err = common.VariableFromEnvironment("ARM_SUBSCRIPTION_ID")
	if err != nil {
		return err
	}

	return nil
}
