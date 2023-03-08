package serve

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/controller/host"
	"github.com/hxx258456/pyramidel-chain-baas/controller/member"
	"github.com/hxx258456/pyramidel-chain-baas/controller/orgnizations"
	"github.com/hxx258456/pyramidel-chain-baas/controller/peer"
	"github.com/hxx258456/pyramidel-chain-baas/internal/localconfig"
	"github.com/hxx258456/pyramidel-chain-baas/internal/routers"
	"github.com/hxx258456/pyramidel-chain-baas/internal/version"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ginLogger = zap.L().Named("gin/serve")

func Serve() {
	r := gin.New()
	gin.SetMode(localconfig.Defaultconfig.Serve.Mode)

	//r.Use(whiteListIp(&localconfig.Defaultconfig.Serve), logger.GinLogger(ginLogger, &localconfig.Defaultconfig.Logger), logger.RecoveryWithZap(ginLogger, true))
	r.Use(logger.GinLogger(ginLogger, &localconfig.Defaultconfig.Logger), logger.RecoveryWithZap(ginLogger, true))
	routers.Include(host.Routers, orgnizations.Routers, peer.Routers, member.Routers)
	routers.Init(r)
	zap.L().Info(" ", zap.String("version", version.Version))
	srv := http.Server{
		Addr:    localconfig.Defaultconfig.Serve.Port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ginLogger.Panic("Server error:", zap.Error(err))
		}
	}()

	// 优雅重启
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	ginLogger.Info("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		ginLogger.Fatal("Server Shutdown:", zap.Error(err))
	}
	ginLogger.Info("Server exiting")
}

func whiteListIp(conf *localconfig.Serve) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipWhiteList := map[string]bool{}
		for _, v := range conf.IpWhiteList {
			ipWhiteList[v] = true
		}
		if !ipWhiteList[c.ClientIP()] {
			ginLogger.Debug("invalid IP: " + c.ClientIP())
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"stats": http.StatusForbidden, "msg": "Permission denied"})
		}
		c.Next()
	}
}
