package v1

import "github.com/gin-gonic/gin"

type Core struct{}

func (Core) Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ok.",
		"code":    200,
		"data":    "pong",
	})
}
