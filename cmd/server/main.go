package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/common/server"

	"github.com/ydataai/azure-adapter/pkg/component/usage"
	"github.com/ydataai/azure-adapter/pkg/controller"
	"github.com/ydataai/azure-adapter/pkg/service"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	errChan chan error
)

func main() {
	applicationConfiguration := Configuration{}
	restServiceConfiguration := service.RESTServiceConfiguration{}
	serverConfiguration := server.HTTPServerConfiguration{}
	restControllerConfiguration := config.RESTControllerConfiguration{}
	loggerConfiguration := logging.LoggerConfiguration{}

	if err := config.InitConfigurationVariables([]config.ConfigurationVariables{
		&applicationConfiguration,
		&restServiceConfiguration,
		&serverConfiguration,
		&restControllerConfiguration,
		&loggerConfiguration,
	}); err != nil {
		fmt.Println(fmt.Errorf("could not set configuration variables. Err: %v", err))
		os.Exit(1)
	}

	logger := logging.NewLogger(loggerConfiguration)

	authorizer, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	computeUsageClient := compute.NewUsageClient(applicationConfiguration.SubscriptionID)
	computeUsageClient.Authorizer = authorizer

	usageClient := usage.NewUsageClient(computeUsageClient)

	restService := service.NewRESTService(logger, restServiceConfiguration, usageClient)
	restController := controller.NewRESTController(logger, restService, restControllerConfiguration)

	serverCtx := context.Background()
	httpServer := server.NewServer(logger, serverConfiguration)
	restController.Boot(httpServer)
	httpServer.Run(serverCtx)

	for err := range errChan {
		logger.Error(err)
	}
}
