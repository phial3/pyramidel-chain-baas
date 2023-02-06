package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/internal/localconfig"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"go.uber.org/zap"
)

var ginLogger = zap.L().Named("gin/serve")

func SetUpRouter() {
	r := gin.New()

	r.Use(logger.GinzapWithConfig(ginLogger, &localconfig.Defaultconfig.Logger), logger.RecoveryWithZap(ginLogger, true))

	r.GET("/", func(ctx *gin.Context) {

		ginLogger.Panic("An unexpected error happen!!@#################")

	})
	err := r.Run(localconfig.Defaultconfig.Serve.Port)
	if err != nil {
		return
	}

}
