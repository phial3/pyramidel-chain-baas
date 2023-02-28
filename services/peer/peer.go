package peer

import (
	"github.com/docker/docker/api/types"
	"github.com/hxx258456/pyramidel-chain-baas/model"
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
	return nil, nil
}

func (s *peerService) i() {
}

func (s *peerService) Create() {
}
