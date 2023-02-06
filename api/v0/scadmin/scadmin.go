package scadmin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// NewOrgJoin 新组织加入通道
// TODO:
func NewOrgJoin(ctx *gin.Context) {
	ctx.String(http.StatusOK, "新组织加入系统通道")
	return
}
