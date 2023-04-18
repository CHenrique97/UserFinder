package connectDB

import (
	"log"
	"os"

	"github.com/UserFinder/initializers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func init() {
	initializers.LoadEnv()
}

var DB *gorm.DB

func InitConnector() {
	// Replace the values with your PlanetScaleDB credentials
	var err error
	DB, err = gorm.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	log.Println("Successfully connected to PlanetScale!")
}
