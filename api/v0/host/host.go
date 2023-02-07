package host

import (
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/remotessh"
	"go.uber.org/zap"
	"net/http"
)

var hostLogger = zap.L().Named("api/host")

type Host struct {
	Ip      string `json:"ip" binding:"required"`   // ip地址
	Pw      string `json:"pw" binding:"required"`   // root用户密码
	Remark  string `json:"remark"`                  // 备注
	SSHPort uint   `json:"sshPort"`                 // ssh port 为空时默认使用22端口
	Name    string `json:"name" binding:"required"` // 服务器名称
	IsUse   bool   `json:"isUse"`                   // 是否开放使用
}

type Response struct {
	Code int         `json:"code" ` // 响应码
	Msg  string      `json:"msg"`   // 消息
	Data interface{} `json:"data"`  // 数据
}

// Add TODO:添加新主机
func Add(ctx *gin.Context) {
	param := &Host{}

	resp := &Response{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		return
	}

	client, err := remotessh.ConnectToHost(param.Pw, param.Ip, param.SSHPort)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			hostLogger.Error(err.Error())
		}
	}()
	out, err := client.Run("ls -al")
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		return
	}
	ctx.String(http.StatusOK, string(out))
	return
}
