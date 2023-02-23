package organization

import (
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/jsonrpcClient"
	psutilclient "github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/client"
	organizations2 "github.com/hxx258456/pyramidel-chain-baas/pkg/request/organizations"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"github.com/hxx258456/pyramidel-chain-baas/repository/organizations"
	"github.com/hxx258456/pyramidel-chain-baas/services/loadbalance"
)

var orgLogger = logger.Lg.Named("services/organization")

type OrganizationService interface {
	Add(organizations2.Organizations) error
}

type organizationsService struct {
	repo organizations.OrganizationRepo
	lb   loadbalance.LBS
}

func NewOrganizationService() OrganizationService {
	return &organizationsService{}
}

func (s *organizationsService) Add(param organizations2.Organizations) error {
	host := &model.Host{}
	lb, err := host.InitHostLB()

	if err != nil {
		return err
	}
	s.lb = lb
	hostId := s.lb.NextService()
	err = host.QueryById(hostId, host)

	if err != nil {
		return err
	}
	cli, err := jsonrpcClient.ConnetJsonrpc(host.UseIp + ":8082")
	if err != nil {
		return err
	}
	defer cli.Close()
	port, err := psutilclient.CallGetPort(cli)
	if err != nil {
		return err
	}
	org := model.Organization{
		Uscc:           param.OrgUscc,
		CaServerDomain: "ca." + param.OrgUscc + ".com",
		CaServerName:   "ca-" + param.OrgUscc,
		CaHostId:       hostId,
		CaServerPort:   uint(port),
		CaUser:         "admin",
		CaPassword:     param.OrgUscc,
		Domain:         param.OrgUscc + ".pcb.com",
		Status:         0,
	}
	s.repo = &org
	if err := s.repo.Create(param, s.lb); err != nil {
		return err
	}
	return nil
}

func (s *organizationsService) RunContainer() {

}
