package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Os struct{}

func (Os) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"code":    http.StatusOK,
		"data":    gin.H{},
	})
}
