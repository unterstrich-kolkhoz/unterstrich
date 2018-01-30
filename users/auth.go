package users

import (
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// InitializeAuth initializes authentication; it also gets the secret key from
// an environment variable called `ARFX_SECRET_KEY`. Make sure this variable is
// set on production.
func InitializeAuth(db *gorm.DB, r *gin.Engine) func() gin.HandlerFunc {
	secretKey := os.Getenv("ARFX_SECRET_KEY")

	if secretKey == "" {
		secretKey = "DEBUG"
	}

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "user zone",
		Key:        []byte(secretKey),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(username string, pw string, c *gin.Context) (string, bool) {
			var user User
			db.Where("username = ?", username).First(&user)
			if user.ID != 0 && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw)) == nil {
				return username, true
			}

			return username, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TimeFunc: time.Now,
	}

	r.POST("/login", authMiddleware.LoginHandler)

	return authMiddleware.MiddlewareFunc
}
