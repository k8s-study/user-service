package db

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/k8s-study/user-service/models"
)

var db *gorm.DB
var err error

func Init() gin.HandlerFunc {
	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"))

	db, err := gorm.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&models.User{})

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}
