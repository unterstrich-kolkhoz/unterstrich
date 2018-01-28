package users

import (
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/hellerve/artifex/model"
)

type User struct {
	model.Base
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Username  string   `json:"username"`
	artist    bool     `json:"is_artist"`
	curator   bool     `json:"is_curator"`
	admin     bool     `json:"is_admin"`
	staff     bool     `json:"is_staff"`
	Address   *Address `json:"address"`
	Social    *Social  `json:"social"`
}
type Address struct {
	model.Base
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	City  string `json:"city"`
	State string `json:"state"`
}
type Social struct {
	model.Base
	Github  string `json:"github"`
	Ello    string `json:"ello"`
	Website string `json:"website"`
}

func Initialize(db *gorm.DB, router *gin.Engine) {
	router.GET("/users", endpoint(db, GetUsers))
	router.POST("/users", endpoint(db, CreateUser))
	router.GET("/users/:id", endpoint(db, GetUser))
	router.PUT("/users/:id", endpoint(db, UpdateUser))
	router.DELETE("/users/:id", endpoint(db, DeleteUser))

	db.AutoMigrate(&User{}, &Address{}, &Social{})
}

func endpoint(db *gorm.DB, wrapped func(*gin.Context, *gorm.DB)) func(*gin.Context) {
	return func(c *gin.Context) {
		wrapped(c, db)
	}
}

func GetUsers(c *gin.Context, db *gorm.DB) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "ID must be numerical: ", err.Error())
		return
	}

	var user *User
	db.First(user, id)

	if user == nil {
		c.String(http.StatusNotFound, "Invalid ID: not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context, db *gorm.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	if !db.NewRecord(user) {
		c.String(http.StatusBadRequest, "User already present: ", string(user.ID))
		return
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	user.Password = string(pw)

	db.Create(&user)

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context, db *gorm.DB) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var user *User
	db.First(user, id)

	if user == nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	c.String(http.StatusOK, "")
}

func UpdateUser(c *gin.Context, db *gorm.DB) {
	_, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID: must be numerical")
		return
	}

	var user *User
	if err := c.ShouldBindJSON(user); err != nil {
		c.String(http.StatusBadRequest, "Invalid body: ", err.Error())
		return
	}

	if db.NewRecord(user) {
		c.String(http.StatusNotFound, "Not found")
		return
	}

	db.Save(&user)

	c.JSON(http.StatusOK, user)
}
