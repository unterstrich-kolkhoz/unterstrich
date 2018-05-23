package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unterstrich-kolkhoz/unterstrich/config"
)

type Context struct {
	DB     *gorm.DB
	Config *config.Config
}

func Endpoint(ctx *Context, wrapped func(*gin.Context, *Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		wrapped(c, ctx)
	}
}
