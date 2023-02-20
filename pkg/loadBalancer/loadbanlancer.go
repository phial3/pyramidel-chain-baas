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
