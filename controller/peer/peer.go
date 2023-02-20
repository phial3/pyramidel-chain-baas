package peer

import (
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	peer2 "github.com/hxx258456/pyramidel-chain-baas/services/peer"
	"strconv"
)

type PeerController struct {
	service peer2.PeerService
}

var peer = &PeerController{
	service: peer2.NewPeerService(),
}

func (p *PeerController) Test(ctx *gin.Context) {
	queryId := ctx.Query("id")

	id, err := strconv.Atoi(queryId)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	peer := &model.Peer{}
	if err := peer.GetByID(uint(id)); err != nil {
		response.Error(ctx, err)
		return
	}
	result, err := p.service.List(*peer)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, result, "")
	return
}
