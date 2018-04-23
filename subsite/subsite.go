package subsite

import (
	"log"
	"net/http"
	"os"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/hoisie/mustache"

	"github.com/hellerve/unterstrich/endpoints"
	"github.com/hellerve/unterstrich/users"
)

// Initialize initializes the subsite router context
func Initialize(ctx *endpoints.Context, router *gin.Engine, auth func() gin.HandlerFunc) {
	g := router.Group("/subsite")
	g.Use(auth())
	{
		g.POST("/", endpoints.Endpoint(ctx, UpdateSubsite))
	}
}

func uploadFiles(username string, files []string) {
	// create S3 bucket for user if necessary
	for _, file := range files {
		// upload file to S3 bucket
		err := os.Remove(file)
		if err != nil {
			log.Println("Error during subsite creation (while deleting file", file,
				"): ", err)
		}
	}
}

func processUpdate(ctx *endpoints.Context, username string) {
	var user users.User
	ctx.DB.Where("username = ?", username).First(&user)
	var files []string

	var artworks []users.Artwork
	ctx.DB.Model(&user).Where("selected = ?", true).Related(&artworks)
	data := mustache.RenderFile(ctx.Config.TemplateDir+"/subsite.html", map[string]interface{}{"user": user, "artworks": artworks})

	f, err := os.Create(username + ".html")
	if err != nil {
		log.Println("Error during subsite creation: ", err)
		return
	}
	defer f.Close()
	files = append(files, f.Name())

	_, err = f.Write([]byte(data))
	if err != nil {
		log.Println("Error during subsite creation: ", err)
		return
	}

	data = mustache.RenderFile(ctx.Config.TemplateDir+"/subsite_about.html", map[string]interface{}{"user": user})

	f, err = os.Create(username + "_about.html")
	if err != nil {
		log.Println("Error during subsite creation: ", err)
		return
	}
	defer f.Close()
	files = append(files, f.Name())

	_, err = f.Write([]byte(data))
	if err != nil {
		log.Println("Error during subsite creation: ", err)
		return
	}

	for _, artwork := range artworks {
		data = mustache.RenderFile(ctx.Config.TemplateDir+"/subsite_artwork.html", map[string]interface{}{"user": user, "artwork": artwork})

		f, err = os.Create(username + "_" + artwork.Slug() + ".html")
		files = append(files, f.Name())
		if err != nil {
			log.Println("Error during subsite creation: ", err)
			return
		}
		defer f.Close()

		_, err = f.Write([]byte(data))
		if err != nil {
			log.Println("Error during subsite creation: ", err)
			return
		}
	}
	uploadFiles(username, files)
}

// UpdateSubsite updates the user subsite
func UpdateSubsite(c *gin.Context, ctx *endpoints.Context) {
	claims := jwt.ExtractClaims(c)
	go processUpdate(ctx, claims["id"].(string))

	c.String(http.StatusOK, "")
}
