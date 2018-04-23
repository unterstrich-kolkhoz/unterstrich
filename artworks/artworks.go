package artworks

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	"github.com/hellerve/unterstrich/endpoints"
	"github.com/hellerve/unterstrich/users"
)

// ArtworkJSON is the JSON artwork model
type ArtworkJSON struct {
	Type        string      `json:"type" binding:"required"`
	URL         string      `json:"url" gorm:"unique"`
	Thumbnail   string      `json:"thumbnail"`
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Views       int         `json:"views"`
	UserID      *users.User `json:"owner"`
	Price       float64     `json:"price"`
}

// Initialize installs all endpoints needed for artworks
func Initialize(ctx *endpoints.Context, router *gin.Engine, auth func() gin.HandlerFunc) {
	g := router.Group("/artworks")
	g.Use(auth())
	{
		g.GET("/", endpoints.Endpoint(ctx, GetArtworks))
		g.POST("/", endpoints.Endpoint(ctx, CreateArtwork))
		g.GET("/:id", endpoints.Endpoint(ctx, GetArtwork))
		g.PUT("/:id", endpoints.Endpoint(ctx, UpdateArtwork))
		g.DELETE("/:id", endpoints.Endpoint(ctx, DeleteArtwork))
		g.POST("/:id/upload", endpoints.Endpoint(ctx, UploadArtwork))
	}
}

// GetArtworks gets all artworks
func GetArtworks(c *gin.Context, ctx *endpoints.Context) {
	claims := jwt.ExtractClaims(c)
	var user users.User
	ctx.DB.Where("username = ?", claims["id"]).First(&user)

	var artworks []users.Artwork
	ctx.DB.Model(&user).Related(&artworks)
	c.JSON(http.StatusOK, artworks)
}

// GetArtwork gets all artworks; adds a view if it’s not the owner
func GetArtwork(c *gin.Context, ctx *endpoints.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "ID must be numerical: ", err.Error())
		return
	}

	var artwork users.Artwork
	ctx.DB.First(&artwork, id)

	if artwork.ID == 0 {
		c.String(http.StatusNotFound, "Invalid ID: not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	ctx.DB.Where("username = ?", claims["id"]).First(&user)

	if user.ID != artwork.UserID {
		artwork.Views++
		ctx.DB.Save(&artwork)
	}

	c.JSON(http.StatusOK, artwork)
}

func createThumbnail(art *users.Artwork, ctx *endpoints.Context) {
	var marshalled []byte
	var err error
	switch art.Type {
	case "image":
		marshalled, err = json.Marshal(gin.H{
			"width":       300,
			"compression": 80,
			"format":      "png",
			"url":         art.URL,
		})
	case "video":
		marshalled, err = json.Marshal(gin.H{
			"width": 300,
			"url":   art.URL,
		})
	default:
		return
	}

	if err != nil {
		log.Println("Error while generating thumbnail for artwork ", art.ID,
			": ", err)
		return
	}

	buf := bytes.NewBuffer(marshalled)
	var resp *http.Response
	switch art.Type {
	case "image":
		resp, err = http.Post(ctx.Config.ImgThumbnailer, "application/json", buf)
	case "video":
		resp, err = http.Post(ctx.Config.VidThumbnailer, "application/json", buf)
	case "default":
		return
	}

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
	ctx.DB.Save(art)
}

// CreateArtwork creates an artwork
func CreateArtwork(c *gin.Context, ctx *endpoints.Context) {
	var artjson ArtworkJSON
	if err := c.ShouldBindJSON(&artjson); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	var art users.Artwork

	art.Type = artjson.Type
	art.Name = artjson.Name
	art.Description = artjson.Description
	art.Price = artjson.Price

	claims := jwt.ExtractClaims(c)
	var user users.User
	ctx.DB.Where("username = ?", claims["id"]).First(&user)

	if &user == nil {
		c.String(http.StatusUnauthorized, "")
		return
	}

	art.UserID = user.ID

	if !ctx.DB.NewRecord(art) {
		c.String(http.StatusBadRequest, "Artwork already present: ", string(art.ID))
		return
	}

	ctx.DB.Create(&art)

	c.JSON(http.StatusOK, art)
}

func processUploadArtwork(file multipart.File, art *users.Artwork, ctx *endpoints.Context) {
	var b bytes.Buffer
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("Error while uploading artwork ", art.ID,
				": ", err)
		}
	}()

	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("upload", "upload")

	if err != nil {
		log.Println("Error while uploading artwork ", art.ID,
			": ", err)
		return
	}

	if _, err = io.Copy(fw, file); err != nil {
		log.Println("Error while uploading artwork ", art.ID,
			": ", err)
		return
	}

	defer func() {
		err := w.Close()
		if err != nil {
			log.Println("Error while uploading artwork ", art.ID,
				": ", err)
		}
	}()

	go func() {
		req, err := http.NewRequest("POST", ctx.Config.Uploader, &b)
		if err != nil {
			log.Println("Error while uploading artwork ", art.ID,
				": ", err)
			return
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Println("Error while uploading artwork ", art.ID,
				": ", err)
			return
		}

		type Response struct {
			URL string `json:"url"`
		}

		var content Response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while uploading artwork ", art.ID,
				": ", err)
			return
		}

		err = json.Unmarshal(body, &content)
		if err != nil {
			log.Println("Error while uploading artwork ", art.ID,
				": ", err)
			return
		}

		art.URL = content.URL
		ctx.DB.Save(art)

		go createThumbnail(art, ctx)
	}()
}

// UploadArtwork uploads an actual artwork to S3
func UploadArtwork(c *gin.Context, ctx *endpoints.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var art users.Artwork
	ctx.DB.First(&art, id)

	if &art == nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	ctx.DB.Where("username = ?", claims["id"]).First(&user)

	if user.ID != art.UserID {
		c.String(http.StatusForbidden, "Cannot alter foreign artwork")
		return
	}

	file, _, err := c.Request.FormFile("upload")

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid file")
		return
	}

	processUploadArtwork(file, &art, ctx)

	c.String(http.StatusOK, "")
}

// DeleteArtwork deletes an artwork
func DeleteArtwork(c *gin.Context, ctx *endpoints.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var art *users.Artwork
	ctx.DB.First(art, id)

	if art.ID == 0 {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	ctx.DB.Where("username = ?", claims["id"]).First(&user)

	if user.ID != art.UserID {
		c.String(http.StatusForbidden, "Cannot alter foreign artwork")
		return
	}

	ctx.DB.Delete(art)

	c.String(http.StatusOK, "")
}

// UpdateArtwork updates an artwork
func UpdateArtwork(c *gin.Context, ctx *endpoints.Context) {
	_, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var art *users.Artwork
	if err := c.ShouldBindJSON(art); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	if ctx.DB.NewRecord(art) {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var user users.User
	ctx.DB.Where("username = ?", claims["id"]).First(&user)

	if user.ID != art.UserID {
		c.String(http.StatusForbidden, "Cannot alter foreign artwork")
		return
	}

	ctx.DB.Save(&art)

	c.JSON(http.StatusOK, art)
}
