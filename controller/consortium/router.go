package consortium

import (
	"github.com/gin-gonic/gin"
)

func Routers(e *gin.RouterGroup) {
	routerG := e.Group("/consortium")
	routerG.POST("/new", consortium.New)
	routerG.POST("/update", consortium.Update)
	routerG.GET("/quit", consortium.Quit)
}
