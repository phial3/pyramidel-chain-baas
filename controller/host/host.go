package host

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	shost "github.com/hxx258456/pyramidel-chain-baas/services/host"
	"net/http"
	"strconv"
)

var hostLogger = logger.Lg.Named("controller/host")
var host = Host{
	service: shost.NewHostService(),
}

type Host struct {
	service shost.HostService
}

func (s *Host) Verify(ctx *gin.Context) {
	param := &model.Host{}

	resp := &response.Response{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}
	ip, covRTT, err := s.service.Verify(param)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}
	resp.Code = 1
	format := "use ip[%s] cover RTT: %dms"
	msg := fmt.Sprintf(format, ip, covRTT)
	resp.Msg = msg
	ctx.JSON(http.StatusOK, resp)
	return
}

// Add 异步添加新主机
func (s *Host) Add(ctx *gin.Context) {
	param := &model.Host{}

	resp := &response.Response{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}

	if err := s.service.Add(param); err != nil {
		hostLogger.Error(err.Error())
		resp.Code = 0
		resp.Msg = err.Error()
		ctx.JSON(http.StatusOK, resp)
		ctx.Error(err)
		return
	}
	resp.Code = 1
	resp.Msg = "Init host process working!!!"
	ctx.JSON(http.StatusOK, resp)
	return
}

func (s *Host) List(ctx *gin.Context) {
	resp := response.Response{}
	host := new(model.Host)
	result, err := s.service.List(host)
	if err != nil {
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
func (s *Host) GetResource(ctx *gin.Context) {
	queryId := ctx.Query("id")
	resp := response.Response{}
	id, err := strconv.Atoi(queryId)
	if err != nil {
		resp.Code = 0
		resp.Msg = err.Error()
		resp.Data = nil
		ctx.JSON(http.StatusOK, resp)
		return
	}
	host := &model.Host{}
	checkInfo, err := s.service.GetResource(id, host)
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
