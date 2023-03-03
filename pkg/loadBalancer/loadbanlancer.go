package loadBalancer

import (
	"sync"
)

type Service struct {
	ID             uint
	Weight         int
	CurrentConns   int
	MaxConnections int
	mu             sync.Mutex
}

// TODO:添加权重动态更新算法
// TODO:lb均衡器全局化
// TODO:添加连接数最大限制
type LoadBalancer struct {
	Services []Service
}

func (lb *LoadBalancer) NextService() uint {
	var leastConnService *Service
	for i := range lb.Services {
		service := &lb.Services[i]
		if leastConnService == nil || service.CurrentConns < leastConnService.CurrentConns {
			leastConnService = service
		}
	}

	if leastConnService == nil {
		return 0
	}

	leastConnService.mu.Lock()
	defer leastConnService.mu.Unlock()
	if leastConnService.CurrentConns >= leastConnService.MaxConnections {
		return 0
	}

	leastConnService.CurrentConns++
	return leastConnService.ID
}
