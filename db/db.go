package db

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/k8s-study/user-service/models"
)

var db *gorm.DB
var err error

func Init() gin.HandlerFunc {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=users password=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}


	db.AutoMigrate(&models.User{})

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}
