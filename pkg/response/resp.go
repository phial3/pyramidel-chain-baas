package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code" ` // 响应码
	Msg  string      `json:"msg"`   // 消息
	Data interface{} `json:"data"`  // 数据
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, &Response{
		Code: 0,
		Msg:  err.Error(),
		Data: nil,
	})
	if err := ctx.Error(err); err != nil {

	}
	return
}
func Success(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, &Response{
		Code: 1,
		Msg:  msg,
		Data: data,
	})
	return
}
