package db

import (
	"fmt"
  "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
  gorm.Model
  Email        string  `gorm:"type:varchar(100);unique_index"`
  Password     string  `gorm:"size:255"`
}

var db *gorm.DB
var err error

func Init() {
    db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=users password=postgres sslmode=disable")
    if err != nil {
        fmt.Println(err)
    }

    db.AutoMigrate(&User{})
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
    db.Close()
}
