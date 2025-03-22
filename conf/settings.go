package conf

import "github.com/gin-gonic/gin"

var MIDDLEWARE = []gin.HandlerFunc{
	gin.Recovery(),
	gin.Logger(),
}

var ROUTERS = []func(group *gin.RouterGroup){}
