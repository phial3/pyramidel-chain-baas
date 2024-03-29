package consortium

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/configtxgen"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var consortium = newConsortiumController()

type consortiumController struct {
}

func newConsortiumController() *consortiumController {
	return &consortiumController{}
}

//
func (c *consortiumController) New(ctx *gin.Context) {
	// 创建联盟
	params := &NewReq{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		response.Error(ctx, err)
		return
	}

	temCfg, err := configtxgen.ReadTemplateFile()
	if err != nil {
		response.Error(ctx, err)
		return
	}

	var organizations []*configtxgen.Organization
	var ordererAddress []string
	var consenters []*configtxgen.Consenter
	for k, v := range params.Group {
		var ordererEndpoints []string

		for _, o := range v.Orderer {
			ordererNodePort, err := strconv.Atoi(o.NodePort)
			if err != nil {
				response.Error(ctx, err)
				return
			}
			ordererEndpoints = append(ordererEndpoints, o.Domain+":"+o.NodePort)
			consenters = append(consenters, &configtxgen.Consenter{
				Host:          o.Domain,
				Port:          uint32(ordererNodePort),
				ClientTlsCert: "/root/txhyjuicefs/organizations/" + k + "/" + o.Domain + "/tls/server.crt",
				ServerTlsCert: "/root/txhyjuicefs/organizations/" + k + "/" + o.Domain + "/tls/server.crt",
			})
		}

		ordererAddress = append(ordererAddress, ordererEndpoints...)
		var anchorPeers []*configtxgen.AnchorPeer
		for _, p := range v.Peer {
			peeerNodePort, err := strconv.Atoi(p.NodePort)
			if err != nil {
				response.Error(ctx, err)
				return
			}
			anchorPeers = append(anchorPeers, &configtxgen.AnchorPeer{
				Port: peeerNodePort,
				Host: p.Domain,
			})
		}
		organization := &configtxgen.Organization{
			Name:             k + "MSP",
			ID:               k + "MSP",
			MSPDir:           "/root/txhyjuicefs/organizations/" + k + "/msp",
			Policies:         temCfg.Organizations[0].Policies,
			AnchorPeers:      anchorPeers,
			OrdererEndpoints: ordererEndpoints,
		}
		organizations = append(organizations, organization)
	}
	temCfg.Organizations = organizations

	temCfg.Orderer.OrdererType = configtxgen.EtcdRaft
	temCfg.Orderer.Addresses = ordererAddress

	temCfg.Orderer.EtcdRaft.Consenters = consenters
	temCfg.Orderer.Organizations = organizations

	profiles := map[string]*configtxgen.Profile{}
	temCfg.Application.Organizations = organizations
	temCfg.Application.Capabilities = map[string]bool{"V2_0": true}
	profiles["ApplicationChannel"] = &configtxgen.Profile{
		Consortium:  "SampleConsortium",
		Application: temCfg.Application,
	}
	cfg := configtxgen.TopLevel{
		Profiles:      profiles,
		Organizations: organizations,
	}

	bytes, err := yaml.Marshal(&cfg)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	if err := os.MkdirAll("/root/txhyjuicefs/configtxs/"+params.Initiator+"/", os.ModeDir); err != nil {
		response.Error(ctx, err)
		return
	}
	channelId := fmt.Sprintf("%s_%d", params.Initiator, time.Now().UnixNano())
	filename := fmt.Sprintf("/root/txhyjuicefs/configtxs/%s/%s.yaml", params.Initiator, channelId)
	if err := ioutil.WriteFile(filename, bytes, 0644); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{
		"ChannelId": channelId,
		"configtx":  cfg,
	}, "")
	return
}

func (c *consortiumController) Update(ctx *gin.Context) {
	// 联盟配置更新
}

func (c *consortiumController) Quit(ctx *gin.Context) {
	// 联盟节点退出
}
