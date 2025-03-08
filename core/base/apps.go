package base

import (
	"github.com/gin-gonic/gin"
	baseapiv1 "github.com/gsadism/open-admin/core/base/api/v1"
	basehooks "github.com/gsadism/open-admin/core/base/hooks"
	basemodels "github.com/gsadism/open-admin/core/base/models"
	"github.com/gsadism/open-admin/core/model"
)

func Router(r *gin.RouterGroup) {
	coreAPIV1 := new(baseapiv1.Core)
	v1 := r.Group("/v1")
	{
		v1.GET("ping", coreAPIV1.Ping)
		v1.GET("secret", coreAPIV1.PublicKey)
		v1.GET("verify", coreAPIV1.VerificationCode)
		v1.GET("website", coreAPIV1.WebSite)
	}
}

var Models = []model.IModel{
	basemodels.IRWebsiteSetting{},
}

var Hooks = []func(){
	basehooks.IrWebsiteSettingInit,
}
