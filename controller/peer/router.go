package peer

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {
	scaGroup := e.Group("/peer")
	scaGroup.POST("/test", peer.Test)
}
