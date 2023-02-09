package host

import "github.com/gin-gonic/gin"

func Routers(e *gin.RouterGroup) {
	scaGroup := e.Group("/host")
	scaGroup.POST("/new", Add)
	scaGroup.POST("/verify", Verify)
	scaGroup.GET("/list", List)
}
