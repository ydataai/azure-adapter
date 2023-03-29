// Package metering provides objects to interact with metering API
package metering

import (
	"context"
	"net/http"

	"github.com/ydataai/go-core/pkg/common/config"
	"github.com/ydataai/go-core/pkg/common/server"
	"github.com/ydataai/go-core/pkg/metering"

	"github.com/gin-gonic/gin"
	"github.com/ydataai/go-core/pkg/common/logging"
)

// RESTController defines rest controller
type RESTController struct {
	logger           logging.Logger
	configuration    config.RESTControllerConfiguration
	markeplaceClient metering.Client
}

// NewRESTController initializes rest controller
func NewRESTController(
	logger logging.Logger,
	marketplaceClient metering.Client,
	configuration config.RESTControllerConfiguration,
) RESTController {
	return RESTController{
		logger:           logger,
		configuration:    configuration,
		markeplaceClient: marketplaceClient,
	}
}

// Boot ...
func (r RESTController) Boot(s server.Server) {
	s.Router().POST("/metering/usageEvent", r.usageEvent())
	s.Router().POST("/metering/batchUsageEvent", r.batchUsageEvent())
}

func (r RESTController) usageEvent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tCtx, cancel := context.WithTimeout(ctx, r.configuration.HTTPRequestTimeout)
		defer cancel()

		event := metering.UsageEvent{}
		if err := ctx.ShouldBindJSON(&event); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		r.logger.Infof("got event %+v", event)

		response, err := r.markeplaceClient.CreateUsageEvent(tCtx, event)
		if err != nil {
			r.logger.Errorf("failed with error %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		r.logger.Infof("got response %+v", response)

		ctx.JSON(http.StatusOK, response)
	}
}

func (r RESTController) batchUsageEvent() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tCtx, cancel := context.WithTimeout(ctx, r.configuration.HTTPRequestTimeout)
		defer cancel()

		event := metering.UsageEventBatch{}
		if err := ctx.ShouldBindJSON(&event); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		r.logger.Infof("got event %+v", event)

		response, err := r.markeplaceClient.CreateUsageEventBatch(tCtx, event)
		if err != nil {
			r.logger.Errorf("failed with error %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		r.logger.Infof("got response %+v", response)

		ctx.JSON(http.StatusOK, response)
	}
}
