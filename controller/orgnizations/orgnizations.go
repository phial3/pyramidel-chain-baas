package orgnizations

import (
	"github.com/gin-gonic/gin"
	"github.com/hxx258456/pyramidel-chain-baas/model"
	orgReq "github.com/hxx258456/pyramidel-chain-baas/pkg/request/organizations"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	mlogger "github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
)

var orgLogger = mlogger.Lg.Named("controller/orgnizations")

type Organization struct {
}

func (o Organization) New(ctx *gin.Context) {
	param := orgReq.Organizations{}
	if err := ctx.BindJSON(&param); err != nil {
		response.Error(ctx, err)
		return
	}

	org := model.Organization{
		Uscc:           param.OrgUscc,
		CaServerDomain: "ca." + param.OrgUscc + ".com",
		CaServerName:   "ca-" + param.OrgUscc,
		CaHostId:       model.Host{},
		CaServerPort:   10054,
		CaUser:         "admin",
		CaPassword:     param.OrgUscc,
		Domain:         param.OrgUscc + ".example.com",
	}
	if err := org.Create(); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, nil, "创建成功")
	return
}
