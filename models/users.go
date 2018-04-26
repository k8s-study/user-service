package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password string `gorm:"size:255" json:"-"`
	KongId   string `gorm:"size:255" json:"-"`
}
