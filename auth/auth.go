package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	// Replace with your Google OAuth2 client ID and secret
	googleOauthConfig = &oauth2.Config{
		ClientID:     "YOUR_GOOGLE_CLIENT_ID",
		ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	}
)

func handleHome(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to my Go application!")
}

func handleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	err, _ := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to exchange token")
		return
	}

	// Use the access token to make authenticated requests to Google APIs
	// ...

	c.String(http.StatusOK, "You have successfully logged in with Google!")
}
