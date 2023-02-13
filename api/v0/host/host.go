package host

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/jsonrpcClient"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/remotessh"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"github.com/melbahja/goph"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var hostLogger = logger.Lg.Named("api/host")

type Response struct {
	Code int         `json:"code" ` // 响应码
	Msg  string      `json:"msg"`   // 消息
	Data interface{} `json:"data"`  // 数据
}

// Verify 验证主机
func Verify(ctx *gin.Context) {
	param := &model.Host{}

	resp := &Response{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}
	var ip string
	var covRtt int64
	iipRtt := remotessh.Ping(param.IntranetIp)
	if iipRtt <= 0 {
		pipRtt := remotessh.Ping(param.PublicIp)
		if pipRtt <= 0 {
			resp.Code = 0
			resp.Msg = "invalid publicIp and internalIp"
			ctx.JSON(http.StatusOK, resp)
			ctx.Error(errors.New("invalid publicIp and internalIp"))
			return
		} else {
			ip = param.PublicIp
			covRtt = pipRtt
		}
	} else {
		ip = param.IntranetIp
		covRtt = iipRtt
	}
	client, err := remotessh.ConnectToHost(param.Username, param.Pw, ip, param.SSHPort)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}
	defer func(client *goph.Client) {
		err := client.Close()
		if err != nil {
			hostLogger.Error(err.Error())
		}
	}(client)
	resp.Code = 1
	format := "use ip[%s] cover RTT: %dms"
	msg := fmt.Sprintf(format, ip, covRtt)
	resp.Msg = msg
	ctx.JSON(http.StatusOK, resp)
	return
}

// Add 异步添加新主机
func Add(ctx *gin.Context) {
	param := &model.Host{}

	resp := &Response{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}

	iipRtt := remotessh.Ping(param.IntranetIp)
	if iipRtt <= 0 {
		pipRtt := remotessh.Ping(param.PublicIp)
		if pipRtt <= 0 {
			resp.Code = 0
			resp.Msg = "invalid publicIp and internalIp"
			ctx.JSON(http.StatusOK, resp)
			ctx.Error(errors.New("invalid publicIp and internalIp"))
			return
		} else {
			param.UseIp = param.PublicIp
		}
	} else {
		param.UseIp = param.IntranetIp
	}
	client, err := remotessh.ConnectToHost(param.Username, param.Pw, param.UseIp, param.SSHPort)
	if err != nil {
		hostLogger.Error(err.Error())
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}

	if err := param.Create(); err != nil {
		hostLogger.Error(err.Error())
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}
	go func(client *goph.Client, host *model.Host) {
		defer func() {
			if err := client.Close(); err != nil {
				hostLogger.Error(err.Error())
			}
		}()
		// 异步进行
		_, ok, err := remotessh.InitHost(client)
		if !ok || err != nil {
			if err := host.Update(model.Host{Status: 0}); err != nil {
				hostLogger.Error(err.Error())
				return
			}
		}
		if err := host.Update(model.Host{Status: 1}); err != nil {
			hostLogger.Error(err.Error())
			return
		}
	}(client, param)
	resp.Code = 1
	resp.Msg = "Init host process working!!!"
	ctx.JSON(http.StatusOK, resp)
	return
}

func List(ctx *gin.Context) {
	resp := Response{}
	host := new(model.Host)
	var result []model.Host
	if err := host.QueryAll(&result); err != nil {
		hostLogger.Error(err.Error())
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}
	resp.Code = 1
	resp.Data = result
	ctx.JSON(http.StatusOK, resp)
	return
}

//GetResource 获取服务器资源实时信息
func GetResource(ctx *gin.Context) {
	queryId := ctx.Query("id")
	resp := Response{}
	id, err := strconv.Atoi(queryId)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		resp.Data = nil
		ctx.JSON(http.StatusOK, resp)
		return
	}
	host := &model.Host{}
	if err := host.QueryById(id, host); err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		resp.Data = nil
		ctx.JSON(http.StatusOK, resp)
		return
	}
	hostLogger.Info(" :::::::::::::::", zap.Any("host", host))
	address := fmt.Sprintf("%s:%d", host.PublicIp, 8082)
	cli, err := jsonrpcClient.ConnetJsonrpc(address)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		resp.Data = nil
		ctx.JSON(http.StatusOK, resp)
		return
	}
	defer func() {
		if err := cli.Close(); err != nil {
			hostLogger.Error("关闭grpc客户端时发生错误", zap.Error(err))
		}
	}()
	checkInfo, err := jsonrpcClient.CallPsutil(cli)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		resp.Data = nil
		ctx.JSON(http.StatusOK, resp)
		return
	}
	resp.Code = 1
	resp.Data = checkInfo
	ctx.JSON(http.StatusOK, resp)
	return
}
