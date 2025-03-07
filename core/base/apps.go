package base

import (
	"github.com/gin-gonic/gin"
	baseapiv1 "github.com/gsadism/open-admin/core/base/api/v1"
)

func Router(r *gin.RouterGroup) {
	coreAPIV1 := new(baseapiv1.Core)
	v1 := r.Group("/v1")
	{
		v1.GET("ping", coreAPIV1.Ping)
	}
}
