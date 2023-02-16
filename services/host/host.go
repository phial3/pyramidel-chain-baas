package host

import (
	"errors"
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/jsonrpcClient"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
	psutilclient "github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/client"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/remotessh"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"github.com/hxx258456/pyramidel-chain-baas/repository/host"
	"github.com/melbahja/goph"
	"go.uber.org/zap"
)

var _ HostService = (*hostService)(nil)
var hostLogger = logger.Lg.Named("services/host")

type HostService interface {
	i()
	Verify(*model.Host) (string, int64, error)
	Add(*model.Host) error
	List(*model.Host) ([]model.Host, error)
	GetResource(*model.Host) (check.HostInfo, error)
	GetResourceById(int, *model.Host) (check.HostInfo, error)
}

type hostService struct {
	repo host.HostRepo
}

func (s *hostService) i() {}

func (s *hostService) Verify(param *model.Host) (ip string, covRtt int64, err error) {

	iipRtt := remotessh.Ping(param.IntranetIp)
	if iipRtt <= 0 {
		pipRtt := remotessh.Ping(param.PublicIp)
		if pipRtt <= 0 {
			return ip, covRtt, errors.New("invalid publicIp and internalIp")
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
		return ip, covRtt, err
	}
	defer func(client *goph.Client) {
		err := client.Close()
		if err != nil {
			hostLogger.Error(err.Error())
		}
	}(client)
	return ip, covRtt, errors.New("invalid publicIp and internalIp")
}

// Add 异步添加新主机
func (s *hostService) Add(param *model.Host) error {
	iipRtt := remotessh.Ping(param.IntranetIp)
	hostLogger.Debug("", zap.Int64("iipRtt", iipRtt))
	if iipRtt <= 0 {
		pipRtt := remotessh.Ping(param.PublicIp)
		hostLogger.Debug("", zap.Int64("pipRtt", pipRtt))
		if pipRtt <= 0 {
			return errors.New("invalid publicIp and internalIp")
		} else {
			param.UseIp = param.PublicIp
		}
	} else {
		param.UseIp = param.IntranetIp
	}
	client, err := remotessh.ConnectToHost(param.Username, param.Pw, param.UseIp, param.SSHPort)
	if err != nil {
		return err
	}

	s.repo = param
	if err := s.repo.Create(); err != nil {
		return err
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
			if err := s.repo.Update(model.Host{Status: 0}); err != nil {
				hostLogger.Error(err.Error())
				return
			}
		}
		if err := s.repo.Update(model.Host{Status: 1}); err != nil {
			hostLogger.Error(err.Error())
			return
		}
	}(client, param)
	return nil
}

func (s *hostService) List(param *model.Host) ([]model.Host, error) {
	s.repo = param
	var result []model.Host
	if err := s.repo.QueryAll(&result); err != nil {
		return nil, err
	}
	return result, nil
}

//GetResource 获取服务器资源实时信息
func (s *hostService) GetResource(param *model.Host) (info check.HostInfo, err error) {
	s.repo = param
	address := fmt.Sprintf("%s:%d", param.UseIp, 8082)
	cli, err := jsonrpcClient.ConnetJsonrpc(address)
	if err != nil {
		return info, err
	}
	defer func() {
		if err := cli.Close(); err != nil {
			hostLogger.Error("关闭grpc客户端时发生错误", zap.Error(err))
		}
	}()
	checkInfo, err := psutilclient.CallPsutil(cli)
	if err != nil {
		return info, err
	}
	return checkInfo, nil
}

//GetResourceById 根据id获取服务器资源实时信息
func (s *hostService) GetResourceById(id int, param *model.Host) (info check.HostInfo, err error) {
	s.repo = param
	if err := s.repo.QueryById(id, param); err != nil {
		return info, err
	}
	address := fmt.Sprintf("%s:%d", param.UseIp, 8082)
	cli, err := jsonrpcClient.ConnetJsonrpc(address)
	if err != nil {
		return info, err
	}
	defer func() {
		if err := cli.Close(); err != nil {
			hostLogger.Error("关闭grpc客户端时发生错误", zap.Error(err))
		}
	}()
	checkInfo, err := psutilclient.CallPsutil(cli)
	if err != nil {
		return info, err
	}
	return checkInfo, nil
}

var _ HostService = &hostService{}

func NewHostService() HostService {
	return &hostService{
		repo: &model.Host{},
	}
}
