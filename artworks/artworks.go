package artworks

import (
	"net/http"
	"strconv"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hellerve/artifex/endpoints"
	"github.com/hellerve/artifex/model"
	"github.com/hellerve/artifex/users"
)

type Artwork struct {
	model.Base
	Type   string       `json:"type"`
	URL    string       `json:"url"`
	Views  int          `json:"views"`
	owner  users.User   `json:"stars"`
	Stars  []users.User `gorm:"many2many:user_languages;" json:"stars"`
	Public bool         `json:"public"`
	Price  float64      `json:"price"`
}

func Initialize(db *gorm.DB, router *gin.Engine, auth func() gin.HandlerFunc) {
	router.GET("/artworks-public", endpoints.Endpoint(db, PublicArtworks))

	g := router.Group("/artworks")
	g.Use(auth())
	{
		g.GET("/", endpoints.Endpoint(db, GetArtworks))
		g.POST("/", endpoints.Endpoint(db, CreateArtwork))
		g.GET("/:id", endpoints.Endpoint(db, GetArtwork))
		g.PUT("/:id", endpoints.Endpoint(db, UpdateArtwork))
		g.DELETE("/:id", endpoints.Endpoint(db, DeleteArtwork))
		g.GET("/:id/star", endpoints.Endpoint(db, StarArtwork))
		g.GET("/:id/unstar", endpoints.Endpoint(db, UnstarArtwork))
	}

	db.AutoMigrate(&Artwork{})
}

func PublicArtworks(c *gin.Context, db *gorm.DB) {
	var artworks []Artwork
	db.Where("public = ?", true).Find(&artworks)
	c.JSON(http.StatusOK, artworks)
}

func GetArtworks(c *gin.Context, db *gorm.DB) {
	var artworks []Artwork
	db.Find(&artworks)
	c.JSON(http.StatusOK, artworks)
}

func GetArtwork(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "ID must be numerical: ", err.Error())
		return
	}

	var artwork *Artwork
	db.First(artwork, id)

	if artwork == nil {
		c.String(http.StatusNotFound, "Invalid ID: not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	db.Where("username = ?", claims["id"]).First(&user)

	if user.ID != artwork.owner.ID {
		artwork.Views += 1
		db.Save(&artwork)
	}

	c.JSON(http.StatusOK, artwork)
}

func CreateArtwork(c *gin.Context, db *gorm.DB) {
	var art Artwork
	if err := c.ShouldBindJSON(&art); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	if !db.NewRecord(art) {
		c.String(http.StatusBadRequest, "Artwork already present: ", string(art.ID))
		return
	}

	db.Create(&art)

	c.JSON(http.StatusOK, art)
}

func DeleteArtwork(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var art *Artwork
	db.First(art, id)

	if art == nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	db.Delete(art)

	c.String(http.StatusOK, "")
}

func UpdateArtwork(c *gin.Context, db *gorm.DB) {
	_, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var art *Artwork
	if err := c.ShouldBindJSON(art); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	if db.NewRecord(art) {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	db.Save(&art)

	c.JSON(http.StatusOK, art)
}

func contains(users []users.User, username string) bool {
	for _, u := range users {
		if u.Username == username {
			return true
		}
	}
	return false
}

func StarArtwork(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var art *Artwork
	db.First(art, id)

	if art == nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	db.Where("username = ?", claims["id"]).First(&user)

	if contains(art.Stars, user.Username) {
		c.String(http.StatusBadRequest, "Artwork is already starred by you")
		return
	}

	db.Model(&art).Association("Stars").Append(user)
	db.Save(&art)

	c.String(http.StatusOK, "")
}

func UnstarArtwork(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var art *Artwork
	db.First(art, id)

	if art == nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	db.Where("username = ?", claims["id"]).First(&user)

	if !contains(art.Stars, user.Username) {
		c.String(http.StatusBadRequest, "Artwork is already starred by you")
		return
	}

	db.Model(&art).Association("Stars").Delete(user)
	db.Save(&art)

	c.String(http.StatusOK, "")
}
