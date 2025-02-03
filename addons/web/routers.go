package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/addons/web/api"
)

var (
	osAPI = new(api.Os)
)

func router(r *gin.RouterGroup) {
	osGroup := r.Group("/")
	{
		osGroup.GET("ping", osAPI.Ping)
		osGroup.GET("rsa", osAPI.PublicKey)
	}
}
