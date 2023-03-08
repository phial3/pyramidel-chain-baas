package member

import "github.com/gin-gonic/gin"

func Routers(e *gin.RouterGroup) {
	routerG := e.Group("/member")
	routerG.POST("/new", member.New)
	routerG.POST("/downloadKS", member.DownloadKeyStore)
}
