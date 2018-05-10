package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/unterstrich-kolkhoz/unterstrich/config"
	"github.com/jinzhu/gorm"
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
