package subsite

import (
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/appleboy/gin-jwt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/hoisie/mustache"

	"github.com/unterstrich-kolkhoz/unterstrich/endpoints"
	"github.com/unterstrich-kolkhoz/unterstrich/users"
)

// Initialize initializes the subsite router context
func Initialize(ctx *endpoints.Context, router *gin.Engine, auth func() gin.HandlerFunc) {
	g := router.Group("/subsite")
	g.Use(auth())
	{
		g.POST("/", endpoints.Endpoint(ctx, UpdateSubsite))
	}
}

//var readGrant = "uri=\"http://acs.amazonaws.com/groups/global/AllUsers\""

func uploadFiles(ctx *endpoints.Context, username string, files []string) {
	usr, err := user.Current()
	if err != nil {
		log.Println("Error during subsite creation, could not get current user")
		return
	}

	dir := usr.HomeDir
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(ctx.Config.Region),
		Credentials: credentials.NewSharedCredentials(dir+"/.aws/credentials", "unterstrich"),
	})

	if err != nil {
		log.Println("Error during subsite creation, could not authenticate to S3")
		return
	}

	s3manage := s3.New(sess)
	bucketName := "unterstrich-" + username
	bucketInput := s3.CreateBucketInput{Bucket: &bucketName}
	// TODO: read grant permissions
	//_, err = s3manage.CreateBucket(bucketInput.SetGrantReadACP(readGrant))
	_, err = s3manage.CreateBucket(&bucketInput)

	if err != nil {
		log.Println("Error during subsite creation, could not create bucket. ", err)
		return
	}

	// create A record if necessary
	for _, file := range files {

		uploader := s3manager.NewUploader(sess)

		f, err := os.Open(file)

		if err != nil {
			log.Println("Error during subsite creation (while opening file", file,
				"): ", err)
			continue
		}

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(file),
			Body:   f,
		})

		if err != nil {
			log.Println("Error during subsite creation (while uploading file", file,
				"): ", err)
			continue
		}

		err = os.Remove(file)
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
	uploadFiles(ctx, username, files)
}

// UpdateSubsite updates the user subsite
func UpdateSubsite(c *gin.Context, ctx *endpoints.Context) {
	claims := jwt.ExtractClaims(c)
	go processUpdate(ctx, claims["id"].(string))

	c.String(http.StatusOK, "")
}
