package models

import "time"

type BaseGormModelFields struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type User struct {
	BaseGormModelFields
	Email    string `gorm:"type:varchar(100);unique_index" json:"email"`
	Password string `gorm:"size:255" json:"-"`
	KongId   string `gorm:"size:255" json:"-"`
}
