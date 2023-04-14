package controllers

import (
	"log"

	connectDB "github.com/UserFinder/connect"
	"github.com/UserFinder/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//Check parses the result of the query
type Check struct {
	Result bool
}

//Migrate migrates the database
func Migrate(c *gin.Context) {
	log.Println("Migrating")
	connectDB.DB.AutoMigrate(&models.User{})
	c.JSON(200, gin.H{
		"message": "Database Migrated",
	})
}

//PostCreate creates a new user
func PostCreate(c *gin.Context) {
	var body models.User
	c.BindJSON(&body)

	post := models.User{
		ID:       uuid.New().String(),
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	}

	var check Check

	err := connectDB.DB.Raw("SELECT EXISTS(SELECT 1 FROM `users` WHERE `email` = ?) as result", body.Email).Scan(&check).Error

	if err != nil {
		// handle error
		c.JSON(400, gin.H{
			"message": "User could not be verified",
		})
		return
	}

	if check.Result {
		c.JSON(400, gin.H{
			"message": "User could not be created",
		})
		return
	}
	// use the `result` variable
	result := connectDB.DB.Create(&post)
	if result.Error != nil {

		c.JSON(400, gin.H{
			"message": "User could not be created",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": &post,
	})
}

//GetUser checks if user exists in database
func GetUser(c *gin.Context) {
	var body models.User
	c.BindJSON(&body)
	user, err := authenticateUser(body.Email, body.Password)

	if err != nil {
		c.JSON(401, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Authentication successful",
		"user":    user.ID,
	})
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
func authenticateUser(email string, password string) (models.User, error) {
	var user models.User
	result := connectDB.DB.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func postDelete(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
