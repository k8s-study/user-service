package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/k8s-study/user-service/models"
)

var db *gorm.DB
var err error

func Init() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=users password=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&models.User{})
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	db.Close()
}
