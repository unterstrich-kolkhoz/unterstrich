package subsite

import (
	"net/http"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hellerve/unterstrich/endpoints"
	"github.com/hellerve/unterstrich/users"
)

// Initialize initializes the subsite router context
func Initialize(db *gorm.DB, router *gin.Engine, auth func() gin.HandlerFunc) {
	g := router.Group("/subsite")
	g.Use(auth())
	{
		g.POST("/", endpoints.Endpoint(db, UpdateSubsite))
	}
}

func processUpdate(db *gorm.DB, username string) {
	var user users.User
	db.Where("username = ?", username).First(&user)

	var artworks []users.Artwork
	db.Model(&user).Where("selected = true").Related(&artworks)
}

// UpdateSubsite updates the user subsite
func UpdateSubsite(c *gin.Context, db *gorm.DB) {
	claims := jwt.ExtractClaims(c)
	go processUpdate(db, claims["id"].(string))

	c.String(http.StatusOK, "")
}
