package host

import "github.com/hxx258456/pyramidel-chain-baas/model"

var _ HostRepo = (*model.Host)(nil)

type HostRepo interface {
	Create() error
	Update(val model.Host) error
	QueryAll(result interface{}) error
	QueryById(uint, interface{}) error
}
