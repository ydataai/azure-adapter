package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/azure-adapter/pkg/common"
	"github.com/ydataai/azure-adapter/pkg/component/marketplace"
	"github.com/ydataai/azure-adapter/pkg/component/usage"
	"github.com/ydataai/azure-adapter/pkg/controller"
	"github.com/ydataai/azure-adapter/pkg/server"
	"github.com/ydataai/azure-adapter/pkg/service"

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

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	computeUsageClient := compute.NewUsageClient(applicationConfiguration.subscriptionID)
	computeUsageClient.Authorizer = authorizer

	marketplaceClient := marketplace.NewMarketplaceClient(logger)
	marketplaceClient.Client.Authorizer = authorizer

	usageClient := usage.NewUsageClient(logger, computeUsageClient)

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
