package netty

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ginResponse(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": code, "msg": msg, "data": data})
}

func ginSuccess(ctx *gin.Context, msg string, data interface{}) {
	ginResponse(ctx, http.StatusOK, msg, data)
}

func ginFail(ctx *gin.Context, code int, msg string) {
	ginResponse(ctx, code, msg, gin.H{})
}
