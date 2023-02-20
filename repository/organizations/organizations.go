package organizations

import (
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/request/organizations"
	"github.com/hxx258456/pyramidel-chain-baas/services/loadbalance"
)

var _ OrganizationRepo = (*model.Organization)(nil)

type OrganizationRepo interface {
	Create(organizations.Organizations, loadbalance.LBS) error
}
