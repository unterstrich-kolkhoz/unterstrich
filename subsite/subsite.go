package subsite

import (
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/appleboy/gin-jwt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/hoisie/mustache"

	"github.com/unterstrich-kolkhoz/unterstrich/config"
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

func createRecordSet(sess *session.Session, conf *config.Config, username string) error {
	svc := route53.New(sess)
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("CREATE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						AliasTarget: &route53.AliasTarget{
							EvaluateTargetHealth: aws.Bool(false),
							HostedZoneId:         aws.String(conf.S3HostedZoneID),
							DNSName:              aws.String(conf.S3URL),
						},
						Name:          aws.String(username + "." + conf.URL),
						Type:          aws.String("A"),
					},
				},
			},
			Comment: aws.String(username + "-creation"),
		},
		HostedZoneId: aws.String(conf.HostedZoneID),
	}
	_, err := svc.ChangeResourceRecordSets(params)
  // A record exists already
  if aerr, ok := err.(awserr.Error); ok && aerr.Code() == route53.ErrCodeInvalidChangeBatch {
    return nil
  }
	return err
}

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
	bucketName := username + "." + ctx.Config.URL
	bucketInput := s3.CreateBucketInput{Bucket: &bucketName}
	_, err = s3manage.CreateBucket(&bucketInput)

	if err != nil && err.(awserr.Error).Code() != s3.ErrCodeBucketAlreadyOwnedByYou {
		log.Println("Error during subsite creation, could not create bucket. ", err)
		return
	}

	// TODO: make me pretty
	policyInput := s3.PutBucketPolicyInput{
		Bucket: &bucketName,
		Policy: aws.String(`{
      "Version": "2012-10-17",
      "Statement": [
          {
              "Sid":"AddPerm",
              "Effect":"Allow",
              "Principal": "*",
              "Action":["s3:GetObject"],
              "Resource": "arn:aws:s3:::` + bucketName + `/*"
          }
      ]
    }`),
	}
	_, err = s3manage.PutBucketPolicy(&policyInput)
	if err != nil {
		log.Println("Error during subsite creation, could not create bucket policy. ", err)
		return
	}

	webconf := s3.WebsiteConfiguration{
		IndexDocument: &s3.IndexDocument{
			Suffix: aws.String("index.html"),
		},
	}
	webinp := s3.PutBucketWebsiteInput{
		Bucket:               &bucketName,
		WebsiteConfiguration: &webconf,
	}
	_, err = s3manage.PutBucketWebsite(&webinp)
	if err != nil {
		log.Println("Error during subsite creation, could not make bucket website. ", err)
		return
	}

	uploader := s3manager.NewUploader(sess)

	// create A record if necessary
	err = createRecordSet(sess, ctx.Config, username)

	if err != nil {
		log.Println("Error during subsite creation, could not create A record. ", err)
		return
	}

	for _, file := range files {
		f, err := os.Open(file)

		if err != nil {
			log.Println("Error during subsite creation (while opening file", file,
				"): ", err)
			continue
		}

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(file),
			ContentType: aws.String("text/html"),
			Body:        f,
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

	f, err := os.Open("templates/logo.png")
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String("logo.png"),
		ContentType: aws.String("image/png"),
		Body:        f,
	})

	if err != nil {
		log.Println("Error during subsite creation (while uploading logo): ", err)
	}
}

func processUpdate(ctx *endpoints.Context, username string) {
	var user users.User
	ctx.DB.Where("username = ?", username).First(&user)
	var files []string

	var artworks []users.Artwork
	ctx.DB.Model(&user).Related(&artworks)
	data := mustache.RenderFile(ctx.Config.TemplateDir+"/subsite.html", map[string]interface{}{"user": user, "artworks": artworks})

	f, err := os.Create("index.html")
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

	f, err = os.Create("about.html")
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
		data = mustache.RenderFile(ctx.Config.TemplateDir+"/subsite_artwork.html", map[string]interface{}{"user": user, "artwork": artwork, "StripeKey": ctx.Config.StripeKey})

		f, err = os.Create(artwork.Slug() + ".html")
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
