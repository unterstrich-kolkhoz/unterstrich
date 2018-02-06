package static

import "github.com/gin-gonic/gin"

// Initialize initializes static file handling
func Initialize(staticdir string, r *gin.Engine) {
	r.Static("/static", staticdir)
	r.StaticFile("/", staticdir+"index.html")
	r.NoRoute(func(g *gin.Context) {
		g.File(staticdir + "index.html")
	})
}
