package users

import (
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"

	"github.com/unterstrich-kolkhoz/unterstrich/endpoints"
	"github.com/unterstrich-kolkhoz/unterstrich/model"
)

// Artwork is the artwork model
type Artwork struct {
	model.Base
	Type        string  `json:"type" binding:"required"`
	URL         string  `json:"url" gorm:"unique"`
	Thumbnail   string  `json:"thumbnail"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Views       int     `json:"views"`
	UserID      uint    `json:"owner"`
	Price       float64 `json:"price"`
	Selected    bool    `json:"selected"`
}

// User is the user model
type User struct {
	model.Base
	Email    string    `json:"email" binding:"required"`
	Password string    `json:"-"`
	Name     string    `json:"name"`
	Username string    `json:"username" binding:"required" gorm:"unique"`
	Admin    bool      `json:"-"`
	Staff    bool      `json:"is_staff"`
	Artworks []Artwork `json:"artworks"`

	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`

	Github  string `json:"github"`
	Ello    string `json:"ello"`
	Website string `json:"website"`
}

// CreationUser is a user model on creation
type CreationUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

// Initialize initializes the URLs for users
func Initialize(ctx *endpoints.Context, router *gin.Engine, auth func() gin.HandlerFunc) {
	router.POST("/users", endpoints.Endpoint(ctx, CreateUser))
	g := router.Group("/users")
	g.Use(auth())
	{
		g.GET("/", endpoints.Endpoint(ctx, GetUsers))
		g.GET("/:id", endpoints.Endpoint(ctx, GetUser))
		g.PUT("/:id", endpoints.Endpoint(ctx, UpdateUser))
		g.DELETE("/:id", endpoints.Endpoint(ctx, DeleteUser))
	}

	g = router.Group("/")
	g.Use(auth())
	{
		g.GET("/me", endpoints.Endpoint(ctx, GetMe))
	}

	ctx.DB.AutoMigrate(&User{}, &Artwork{})
}

// GetUsers gets all users
func GetUsers(c *gin.Context, ctx *endpoints.Context) {
	var users []User
	ctx.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

// GetUser gets a specifc user
func GetUser(c *gin.Context, ctx *endpoints.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "ID must be numerical: ", err.Error())
		return
	}

	var user *User
	ctx.DB.First(user, id)

	if user == nil {
		c.String(http.StatusNotFound, "Invalid ID: not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetMe gets current user
func GetMe(c *gin.Context, ctx *endpoints.Context) {
	claims := jwt.ExtractClaims(c)
	var me User
	ctx.DB.Where("username = ?", claims["id"]).First(&me)

	c.JSON(http.StatusOK, me)
}

func inBlacklist(username string) bool {
	for _, blacklisted := range BlacklistedUsernames {
		if blacklisted == username {
			return true
		}
	}
	return false
}

// CreateUser creates a new user
func CreateUser(c *gin.Context, ctx *endpoints.Context) {
	var jsonuser CreationUser
	if err := c.ShouldBindJSON(&jsonuser); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	var user User

	if inBlacklist(jsonuser.Username) {
		c.String(http.StatusBadRequest, "Username in blacklist: ", jsonuser.Username)
		return
	}

	user.Email = jsonuser.Email
	user.Username = jsonuser.Username
	user.Password = jsonuser.Password
	if !ctx.DB.NewRecord(user) {
		c.String(http.StatusBadRequest, "User already present: ", string(user.ID))
		return
	}

	if user.Staff || user.Admin {
		c.String(http.StatusForbidden, "Cannot create admin user")
		return
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	user.Password = string(pw)

	ctx.DB.Create(&user)

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user (must be the user themselves)
func DeleteUser(c *gin.Context, ctx *endpoints.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var user *User
	ctx.DB.First(user, id)

	if user == nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	claims := jwt.ExtractClaims(c)
	var me User
	ctx.DB.Where("username = ?", claims["id"]).First(&me)

	if user.ID != me.ID {
		c.String(http.StatusForbidden, "Cannot alter foreign user")
		return
	}

	ctx.DB.Delete(&user)

	c.String(http.StatusOK, "")
}

// UpdateUser updates a user (must be the user themselves)
func UpdateUser(c *gin.Context, ctx *endpoints.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	user := &User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	if ctx.DB.NewRecord(user) {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	var other User
	ctx.DB.First(&other, id)

	if (user.Staff && !other.Staff) || (user.Admin && !other.Admin) {
		c.String(http.StatusForbidden, "Cannot make user admin")
		return
	}

	claims := jwt.ExtractClaims(c)
	var me User
	ctx.DB.Where("username = ?", claims["id"]).First(&me)

	if user.ID != me.ID {
		c.String(http.StatusForbidden, "Cannot alter foreign user")
		return
	}

	if len(user.Password) == 0 {
		user.Password = me.Password
	} else {
		pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

		if err != nil {
			c.String(http.StatusInternalServerError, "")
			return
		}

		user.Password = string(pw)
	}

	ctx.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

// Slug generates a slug from the artwork name
func (a Artwork) Slug() string {
	return slug.Make(a.Name)
}

// HasAddress checks whether a user has an address
func (u User) HasAddress() bool {
	return u.Line1 != ""
}

// HasState checks whether a user address contains a state
func (u User) HasState() bool {
	return u.State != ""
}
