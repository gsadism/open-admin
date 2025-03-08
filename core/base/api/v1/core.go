package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/osv"
)

type Core struct{}

func (Core) Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ok.",
		"code":    200,
		"data":    "pong",
	})
}

func (Core) PublicKey(ctx *gin.Context) {
	if key, err := osv.Rsa.PublicKeyWithString(); err != nil {
		ctx.JSON(200, gin.H{
			"message": "fail.",
			"code":    1000,
			"data":    "",
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": "ok.",
			"code":    200,
			"data":    key,
		})
	}
}

func (Core) VerificationCode(ctx *gin.Context) {
	if key, img, err := osv.Image.Number(4); err != nil {
		ctx.JSON(200, gin.H{
			"message": "fail.",
			"code":    1001,
			"data":    "",
		})
	} else {
		ctx.JSON(200, gin.H{
			"message": "ok.",
			"code":    200,
			"data": gin.H{
				"mark": key,
				"img":  img,
			},
		})
	}
}

func (Core) WebSite(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx.JSON(200, gin.H{
				"message": "fail.",
				"code":    1000,
				"data":    gin.H{},
			})
		}
	}()
	db := osv.DB.WithContext(context.WithValue(context.TODO(), "sudo", true))
	var Query struct {
		Title   string `json:"title"`
		Company string `json:"company"`
		ICP     string `json:"icp"`
	}
	db.Table("ir_website_setting").First(&Query)
	ctx.JSON(200, gin.H{
		"message": "ok.",
		"code":    200,
		"data":    Query,
	})
}
