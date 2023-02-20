package loadbalance

import "github.com/hxx258456/pyramidel-chain-baas/pkg/loadBalancer"

// LBS 负载均衡服务
type LBS interface {
	NextService() uint
}

var _ LBS = (*loadBalancer.LoadBalancer)(nil)
