package host

import "github.com/hxx258456/pyramidel-chain-baas/model"

type HostRepo interface {
	Create() error
	Update(val model.Host) error
	QueryAll(result interface{}) error
	QueryById(int, interface{}) error
}
