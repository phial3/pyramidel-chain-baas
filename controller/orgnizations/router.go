package orgnizations

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {
	scaGroup := e.Group("/org")
	scaGroup.POST("/new", org.New)
}
