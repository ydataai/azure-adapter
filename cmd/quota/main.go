// Package main for quota executable
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/logging"
	"github.com/ydataai/go-core/pkg/common/server"

	"github.com/ydataai/azure-adapter/internal/configuration"
	"github.com/ydataai/azure-adapter/internal/usage"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	compute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

var (
	errChan chan error
)

func main() {
	applicationConfiguration := configuration.Application{}
	restServiceConfiguration := usage.RESTServiceConfiguration{}
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

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logger.Fatal(err)
	}

	computeUsageClient, err := compute.NewUsageClient(applicationConfiguration.SubscriptionID, cred, nil)
	if err != nil {
		logger.Fatal(err)
	}
	usageClient := usage.NewClient(computeUsageClient)
	restService := usage.NewRESTService(logger, restServiceConfiguration, usageClient)
	restController := usage.NewRESTController(logger, restService, restControllerConfiguration)

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
