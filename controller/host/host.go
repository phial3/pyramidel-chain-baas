package host

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	shost "github.com/hxx258456/pyramidel-chain-baas/services/host"
	"github.com/jinzhu/copier"
	"strconv"
	"sync"
	"sync/atomic"
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

type res struct {
	result check.HostInfo
	id     int
}

func (s *Host) List(ctx *gin.Context) {
	host := new(model.Host)
	result, err := s.service.List(host)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	resultCh := make(chan res, len(result))
	errCh := make(chan error)
	parentCtx, parentCancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer close(resultCh)
	defer close(errCh)
	count := (int64)(len(result))
	wg.Add(len(result))
	for i, _ := range result {
		go func(parentCtx context.Context, index int, wg *sync.WaitGroup) {
			defer atomic.AddInt64(&count, -1)
			defer wg.Done()
			for {
				select {
				case <-parentCtx.Done():
					return

				default:
					hostinfo, err := s.service.GetResource(&result[index])
					if err != nil {
						errCh <- err
						return
					}
					resultCh <- res{
						result: hostinfo,
						id:     index,
					}
					return
				}
			}
		}(parentCtx, i, &wg)
	}
	for {
		select {
		case err, ok := <-errCh:
			if ok {
				parentCancel()
				wg.Wait()
				response.Error(ctx, err)
				return
			}
		case r, ok := <-resultCh:
			if ok {
				if err := copier.Copy(&result[r.id], &r.result); err != nil {
					parentCancel()
					wg.Wait()
					response.Error(ctx, err)
					return
				}
			}
		default:
			if count <= 0 {
				parentCancel()
				wg.Wait()
				response.Success(ctx, result, "")
				return
			}
		}
	}
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
	checkInfo, err := s.service.GetResourceById(id, host)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, checkInfo, "")
	return
}
