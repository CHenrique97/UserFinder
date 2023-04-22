package main

import (
	"os"

	connectDB "github.com/UserFinder/connect"
	"github.com/UserFinder/controllers"
	"github.com/UserFinder/initializers"
	"github.com/UserFinder/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Main function
func init() {
	initializers.LoadEnv()
	connectDB.InitConnector()
}

func main() {

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your client's URL
	r.Use(cors.New(config))
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.GET("/migrate", controllers.Migrate)
	r.POST("/getUser", controllers.GetUser)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/postUser", controllers.PostCreate)
	r.Run(":" + os.Getenv("port"))

}
