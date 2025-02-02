package t

import "github.com/gin-gonic/gin"

type R = func(*gin.RouterGroup)

type Q struct {
	Routers []R
}
