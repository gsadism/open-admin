package core

import "github.com/gin-gonic/gin"

func ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"msg":  "ok.",
		"code": 200,
		"data": "pong",
	})
}

func routers(r *gin.RouterGroup) {
	r.GET("/ping", ping)
}
