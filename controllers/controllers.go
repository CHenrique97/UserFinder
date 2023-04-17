package controllers

import (
	"fmt"
	"log"

	connectDB "github.com/UserFinder/connect"
	jwtbuilder "github.com/UserFinder/helpers"
	"github.com/UserFinder/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	// Hash the password using bcrypt
	uuid := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password+uuid), bcrypt.DefaultCost)
	if err != nil {
		// handle error
		c.JSON(400, gin.H{
			"message": "User could not be verified",
		})
		return
	}

	post := models.User{
		ID:       uuid,
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword), // Store the hashed password in the database
	}

	var check Check

	err = connectDB.DB.Raw("SELECT EXISTS(SELECT 1 FROM `users` WHERE `email` = ?) as result", body.Email).Scan(&check).Error

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
	token, _ := jwtbuilder.CreateJWTToken(post.ID)

	c.JSON(200, gin.H{

		"message": token,
	})

}

// GetUser gets a user
func GetUser(c *gin.Context) {
	var body models.User
	c.BindJSON(&body)
	//print json from c.BindJSON(&body)

	// Authenticate the user by comparing the hashed password
	user, err := authenticateUser(body.Email, body.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, _ := jwtbuilder.CreateJWTToken(user.ID)

	fmt.Println(token)

	c.JSON(200, gin.H{

		"message": token,
	})

}

func authenticateUser(email string, password string) (models.User, error) {
	var user models.User
	result := connectDB.DB.Where("email = ?", email).First(&user)
	fmt.Println(password)
	if result.Error != nil {
		return user, result.Error
	}

	// Compare the hashed password with the input password

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.ID)); err != nil {

		return user, err
	}

	return user, nil
}

func postDelete(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
