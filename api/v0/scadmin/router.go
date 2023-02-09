package scadmin

import "github.com/gin-gonic/gin"

func Routers(e *gin.RouterGroup) {
	scaGroup := e.Group("/scadmin")
	scaGroup.GET("/newOrg", NewOrgJoin)
}
