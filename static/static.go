package static

import "github.com/gin-gonic/gin"

// Initialize initializes static file handling
func Initialize(staticdir string, r *gin.Engine) {
	r.Static("/static", staticdir)
	r.StaticFile("/", staticdir+"index.html")
}
