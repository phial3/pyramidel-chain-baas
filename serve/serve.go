package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/api/v0/scadmin"
	"github.com/hxx258456/pyramidel-chain-baas/internal/localconfig"
	"github.com/hxx258456/pyramidel-chain-baas/internal/version"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"github.com/hxx258456/pyramidel-chain-baas/routers"
	"go.uber.org/zap"
)

var ginLogger = zap.L().Named("gin/serve")

func Serve() {
	r := gin.New()
	gin.SetMode(localconfig.Defaultconfig.Serve.Mode)

	r.Use(logger.GinzapWithConfig(ginLogger, &localconfig.Defaultconfig.Logger), logger.RecoveryWithZap(ginLogger, true))

	routers.Include(scadmin.Routers)
	routers.Init(r)
	err := r.Run(localconfig.Defaultconfig.Serve.Port)

	if err != nil {
		return
	}
	zap.L().Info(version.Version)
}
