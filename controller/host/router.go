package host

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {
	scaGroup := e.Group("/host")
	scaGroup.POST("/new", host.Add)
	scaGroup.POST("/verify", host.Verify)
	scaGroup.GET("/list", host.List)
	scaGroup.GET("/getResource", host.GetResource)
}
