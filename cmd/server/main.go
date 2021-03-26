package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/azure-quota-provider/pkg/clients"
	"github.com/ydataai/azure-quota-provider/pkg/common"
	"github.com/ydataai/azure-quota-provider/pkg/controller"
	"github.com/ydataai/azure-quota-provider/pkg/server"
	"github.com/ydataai/azure-quota-provider/pkg/service"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"

	"github.com/sirupsen/logrus"
)

func main() {
	applicationConfiguration := configuration{}
	restServiceConfiguration := service.RESTServiceConfiguration{}
	serverConfiguration := server.Configuration{}
	restControllerConfiguration := controller.RESTControllerConfiguration{}

	err := initConfigurationVariables([]common.ConfigurationVariables{
		&applicationConfiguration,
		&restServiceConfiguration,
		&serverConfiguration,
		&restControllerConfiguration,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("could not set configuration variables. Err: %v", err))
		os.Exit(1)
	}

	var logger = logrus.New()
	logger.SetLevel(applicationConfiguration.logLevel)

	computeUsageClient := compute.NewUsageClient(applicationConfiguration.subscriptionID)
	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	computeUsageClient.Authorizer = authorizer

	usageClient := clients.NewUsageClient(logger, computeUsageClient)
	restService := service.NewRESTService(logger, restServiceConfiguration, usageClient)
	restController := controller.NewRESTController(logger, restService, restControllerConfiguration)

	serverCtx := context.Background()
	s := server.NewServer(logger, serverConfiguration)
	restController.Boot(s)
	s.Run(serverCtx)

	for err := range s.ErrCh {
		logger.Error(err)
	}

}

func initConfigurationVariables(configurations []common.ConfigurationVariables) error {
	for _, configuration := range configurations {
		if err := configuration.LoadEnvVars(); err != nil {
			return err
		}
	}

	return nil
}
