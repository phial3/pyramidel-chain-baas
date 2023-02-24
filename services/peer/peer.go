package peer

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/services/container"
)

var _ PeerService = (*peerService)(nil)

type PeerService interface {
	i()
	Create()
	List(peer model.Peer) ([]types.Container, error)
}

type peerService struct {
}

func NewPeerService() PeerService {
	return &peerService{}
}

func (s *peerService) List(peer model.Peer) ([]types.Container, error) {
	ctx := context.Background()
	containerCli := container.NewContainerService(peer.Host.UseIp, peer.Host.PublicIp)
	if err := containerCli.Conn(); err != nil {
		return nil, err
	}
	result, err := containerCli.List(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *peerService) i() {
}

func (s *peerService) Create() {
}
