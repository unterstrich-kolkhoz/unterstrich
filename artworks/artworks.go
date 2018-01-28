package artworks

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hellerve/artifex/endpoints"
	"github.com/hellerve/artifex/users"
	"github.com/hellerve/artifex/model"
)

type Artwork struct {
	model.Base
  Type string `json:"type"`
  URL string `json:"url"`
  Views int `json:"views"`
  Stars []users.User `gorm:"many2many:user_languages;" json:"stars"`
  Public bool `json:"public"`
  Price float64 `json:"price"`
}

func Initialize(db *gorm.DB, router *gin.Engine, auth (func () gin.HandlerFunc)) {
	router.GET("/artworks-public", endpoints.Endpoint(db, PublicArtworks))

	g := router.Group("/artworks")
	g.Use(auth())
	{
		g.GET("/", endpoints.Endpoint(db, GetArtworks))
		g.POST("/", endpoints.Endpoint(db, CreateArtwork))
		g.GET("/:id", endpoints.Endpoint(db, GetArtwork))
		g.PUT("/:id", endpoints.Endpoint(db, UpdateArtwork))
		g.DELETE("/:id", endpoints.Endpoint(db, DeleteArtwork))
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
