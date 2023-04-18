package main

import (
	"os"

	connectDB "github.com/UserFinder/connect"
	"github.com/UserFinder/controllers"
	"github.com/UserFinder/initializers"
	"github.com/UserFinder/middleware"
	"github.com/gin-gonic/gin"
)

// Main function
func init() {
	initializers.LoadEnv()
	connectDB.InitConnector()
}

func main() {

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.GET("/migrate", controllers.Migrate)
	r.GET("/getUser", controllers.GetUser)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/postUser", controllers.PostCreate)
	r.Run(":" + os.Getenv("port"))

}
