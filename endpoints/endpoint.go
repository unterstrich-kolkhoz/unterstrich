package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Endpoint(db *gorm.DB, wrapped func(*gin.Context, *gorm.DB)) func(*gin.Context) {
	return func(c *gin.Context) {
		wrapped(c, db)
	}
}
