package service_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ydataai/azure-adapter/mock"
	"github.com/ydataai/go-core/pkg/common/logging"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/ydataai/azure-adapter/pkg/component/usage"
	"github.com/ydataai/azure-adapter/pkg/service"
)

func TestAvailableGPU(t *testing.T) {
	loggerConfiguration := logging.LoggerConfiguration{}
	if err := loggerConfiguration.LoadFromEnvVars(); err != nil {
		fmt.Println(fmt.Errorf("could not set logging configuration. Err: %v", err))
		os.Exit(1)
	}

	logger := logging.NewLogger(loggerConfiguration)

	t.Run("failure response", func(t *testing.T) {
		errM := errors.New("mock error")

		tt := []struct {
			name        string
			usageClient func(context.Context, *gomock.Controller) usage.UsageClientInterface
			err         error
		}{
			{
				name: "failure on usage client request",
				usageClient: func(ctx context.Context, ctrl *gomock.Controller) usage.UsageClientInterface {
					usageClient := mock.NewMockUsageClientInterface(ctrl)
					usageClient.EXPECT().
						ComputeUsage(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(compute.Usage{}, errM)

					return usageClient
				},
				err: errM,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ctx := context.Background()

				restServiceConfiguration := service.RESTServiceConfiguration{}

				restService := service.NewRESTService(
					logger,
					restServiceConfiguration,
					tc.usageClient(ctx, ctrl),
				)

				_, err := restService.AvailableGPU(ctx)
				if err == nil {
					t.Fatal("should return an error")
				}
			})
		}

	})

	t.Run("successful response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		restServiceConfiguration := service.RESTServiceConfiguration{}

		currentValue := int32(6)
		limit := int64(12)
		usageClient := mock.NewMockUsageClientInterface(ctrl)

		usageClient.EXPECT().
			ComputeUsage(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(compute.Usage{
				CurrentValue: &currentValue,
				Limit:        &limit,
			}, nil)

		restService := service.NewRESTService(
			logger,
			restServiceConfiguration,
			usageClient,
		)

		gpu, err := restService.AvailableGPU(ctx)
		if err != nil {
			t.Fatal("should not return any error")
		}

		if diff := cmp.Diff(gpu, usage.GPU(int64(1))); diff != "" {
			t.Fatalf("should be 1, got %v", gpu)
			t.Fatal(diff)
		}
	})
}
