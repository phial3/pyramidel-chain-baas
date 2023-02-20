package orgnizations

import (
	"github.com/gin-gonic/gin"
	orgReq "github.com/hxx258456/pyramidel-chain-baas/pkg/request/organizations"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/response"
	mlogger "github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"github.com/hxx258456/pyramidel-chain-baas/services/organization"
)

var orgLogger = mlogger.Lg.Named("controller/orgnizations")
var org = Organization{
	service: organization.NewOrganizationService(),
}

type Organization struct {
	service organization.OrganizationService
}

func (o Organization) New(ctx *gin.Context) {
	param := orgReq.Organizations{}
	if err := ctx.BindJSON(&param); err != nil {
		response.Error(ctx, err)
		return
	}
	if err := o.service.Add(param); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, nil, "创建成功")
	return
}
