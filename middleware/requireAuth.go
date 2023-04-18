package middleware

import (
	"net/http"

	connectDB "github.com/UserFinder/connect"
	jwtbuilder "github.com/UserFinder/helpers"
	"github.com/UserFinder/models"
	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware to check if the token is valid
func RequireAuth(c *gin.Context) {
	token, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err = jwtbuilder.VerifyJWTToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	err = connectDB.DB.Where("ID = ?", token).First(&user).Error
	if err != nil || user.ID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user", user)

	c.Next()
}
