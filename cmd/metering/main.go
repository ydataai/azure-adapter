// Package main for metering executable
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/common/server"

	"github.com/ydataai/azure-adapter/internal/configuration"
	"github.com/ydataai/azure-adapter/internal/metering"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var (
	errChan chan error
)

func main() {
	applicationConfiguration := configuration.Application{}
	serverConfiguration := server.HTTPServerConfiguration{}
	restControllerConfiguration := config.RESTControllerConfiguration{}
	loggerConfiguration := logging.LoggerConfiguration{}
	meteringConfiguration := metering.Configuration{}

	if err := config.InitConfigurationVariables([]config.ConfigurationVariables{
		&applicationConfiguration,
		&serverConfiguration,
		&restControllerConfiguration,
		&loggerConfiguration,
		&meteringConfiguration,
	}); err != nil {
		fmt.Println(fmt.Errorf("could not set configuration variables. Err: %v", err))
		os.Exit(1)
	}

	logger := logging.NewLogger(loggerConfiguration)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logger.Fatal(err)
	}

	marketplaceClient, err := metering.NewClient(cred, meteringConfiguration, logger)
	if err != nil {
		logger.Fatal(err)
	}

	restController := metering.NewRESTController(logger, marketplaceClient, restControllerConfiguration)

	serverCtx := context.Background()
	httpServer := server.NewServer(logger, serverConfiguration)
	httpServer.AddHealthz()
	httpServer.AddReadyz(nil)
	restController.Boot(httpServer)
	httpServer.Run(serverCtx)

	for err := range errChan {
		logger.Error(err)
	}
}
