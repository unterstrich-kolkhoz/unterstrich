package artworks

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hellerve/unterstrich/endpoints"
	"github.com/hellerve/unterstrich/model"
	"github.com/hellerve/unterstrich/users"
)

// Artwork is the artwork model
type Artwork struct {
	model.Base
	Type      string       `json:"type" binding:"required"`
	URL       string       `json:"url" binding:"required"`
	Thumbnail string       `json:"thumbnail"`
	Views     int          `json:"views"`
	Owner     *users.User  `json:"owner"`
	Stars     []users.User `gorm:"many2many:user_languages;" json:"stars"`
	Public    bool         `json:"public"`
	Price     float64      `json:"price" binding:"required"`
}

// Initialize installs all endpoints needed for artworks
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

// PublicArtworks gets public artworks
func PublicArtworks(c *gin.Context, db *gorm.DB) {
	var artworks []Artwork
	db.Where("public = ?", true).Find(&artworks)
	c.JSON(http.StatusOK, artworks)
}

// GetArtworks gets all artworks
func GetArtworks(c *gin.Context, db *gorm.DB) {
	var artworks []Artwork
	db.Find(&artworks)
	c.JSON(http.StatusOK, artworks)
}

// GetArtwork gets all artworks; adds a view if it’s not the owner
func GetArtwork(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "ID must be numerical: ", err.Error())
		return
	}

	var artwork Artwork
	db.First(&artwork, id)

	if artwork.ID == 0 {
		c.String(http.StatusNotFound, "Invalid ID: not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	db.Where("username = ?", claims["id"]).First(&user)

	if artwork.Owner == nil || user.ID != artwork.Owner.ID {
		artwork.Views++
		db.Save(&artwork)
	}

	c.JSON(http.StatusOK, artwork)
}

func createThumbnail(art Artwork, db *gorm.DB) {
	marshalled, err := json.Marshal(gin.H{
		"width":       300,
		"compression": 80,
		"format":      "png",
		"url":         art.URL,
	})

	if err != nil {
		log.Println("Error while generating thumbnail for artwork ", art.ID,
			": ", err)
		return
	}

	buf := bytes.NewBuffer(marshalled)
	resp, err := http.Post("http://127.0.0.1:8000/", "application/json", buf)

	if err != nil {
		log.Println("Error while generating thumbnail for artwork ", art.ID,
			": ", err)
		return
	}

	type Response struct {
		URL string `json:"url"`
	}

	var content Response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while generating thumbnail for artwork ", art.ID,
			": ", err)
		return
	}

	err = json.Unmarshal(body, &content)
	if err != nil {
		log.Println("Error while generating thumbnail for artwork ", art.ID,
			": ", err)
		return
	}

	art.Thumbnail = content.URL
	db.Save(&art)
}

// CreateArtwork creates an artwork
func CreateArtwork(c *gin.Context, db *gorm.DB) {
	var art Artwork
	if err := c.ShouldBindJSON(&art); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	db.Where("username = ?", claims["id"]).First(&user)

	art.Owner = &user

	if !db.NewRecord(art) {
		c.String(http.StatusBadRequest, "Artwork already present: ", string(art.ID))
		return
	}

	db.Create(&art)

	go createThumbnail(art, db)

	c.JSON(http.StatusOK, art)
}

// DeleteArtwork deletes an artwork
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

	claims := jwt.ExtractClaims(c)
	var user users.User
	db.Where("username = ?", claims["id"]).First(&user)

	if user.ID != art.Owner.ID {
		c.String(http.StatusForbidden, "Cannot alter foreign artwork")
		return
	}

	db.Delete(art)

	c.String(http.StatusOK, "")
}

// UpdateArtwork updates an artwork
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

	claims := jwt.ExtractClaims(c)
	var user users.User
	db.Where("username = ?", claims["id"]).First(&user)

	if user.ID != art.Owner.ID {
		c.String(http.StatusForbidden, "Cannot alter foreign artwork")
		return
	}

	db.Save(&art)

	c.JSON(http.StatusOK, art)
}

// helper function to test whether a username is in a list of users
func contains(users []users.User, username string) bool {
	for _, u := range users {
		if u.Username == username {
			return true
		}
	}
	return false
}

// StarArtwork stars an artwork (only possible if not already starred)
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

// UnstarArtwork unstars an artwork (only possible if not currently starred)
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
