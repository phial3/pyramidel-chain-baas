package host

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	shost "github.com/hxx258456/pyramidel-chain-baas/services/host"
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
	if err := ctx.ShouldBindJSON(param); err != nil {
		response.Error(ctx, err)
		return
	}
	ip, covRTT, err := s.service.Verify(param)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	format := "use ip[%s] cover RTT: %dms"
	msg := fmt.Sprintf(format, ip, covRTT)
	response.Success(ctx, nil, msg)
	return
}

// Add 异步添加新主机
func (s *Host) Add(ctx *gin.Context) {
	param := &model.Host{}

	if err := ctx.ShouldBindJSON(param); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := s.service.Add(param); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil, "Init host process working!!!")
	return
}

func (s *Host) List(ctx *gin.Context) {
	host := new(model.Host)
	result, err := s.service.List(host)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, result, "")
	return
}

//GetResource 获取服务器资源实时信息
func (s *Host) GetResource(ctx *gin.Context) {
	queryId := ctx.Query("id")

	id, err := strconv.Atoi(queryId)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	host := &model.Host{}
	checkInfo, err := s.service.GetResource(id, host)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, checkInfo, "")
	return
}
