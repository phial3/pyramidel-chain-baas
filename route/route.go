package route

import (
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"net/http"
)

func SetUpRouter() {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	err := r.Run(":8080")
	if err != nil {
		return
	}
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "views/404.html", nil)
	})
}
