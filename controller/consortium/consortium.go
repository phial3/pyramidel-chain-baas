package consortium

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
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

	org := model.Organization{}
	var orgs []model.Organization
	if err := org.FindByUscc(params.Initiator, &orgs); err != nil {
		response.Error(ctx, err)
		return
	}

	if len(orgs) <= 0 {
		response.Error(ctx, errors.New("orgnizations not found!"))
		return
	}

}

func (c *consortiumController) Update(ctx *gin.Context) {
	// 联盟配置更新
}

func (c *consortiumController) Quit(ctx *gin.Context) {
	// 联盟节点退出
}
