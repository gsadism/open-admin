package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/depends/logging"
	"github.com/gsadism/open-admin/osv"
	"github.com/gsadism/open-admin/pkg/crypto/_rsa"
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

func (Os) PublicKey(ctx *gin.Context) {
	if osv.PublicKey == nil {
		logging.Error("rsa public key nil")
		ctx.JSON(http.StatusOK, gin.H{
			"message": "fail.",
			"code":    1000,
			"data":    gin.H{},
		})
	} else {
		if str, err := _rsa.PublicKeyWithString(osv.PublicKey); err != nil {
			logging.Error(err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"message": "fail.",
				"code":    1001,
				"data":    gin.H{},
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "ok.",
				"code":    http.StatusOK,
				"data":    str,
			})
		}
	}
}
