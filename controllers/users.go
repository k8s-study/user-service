package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/k8s-study/user-service/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Customer struct {
	Custom_id string `json:"custom_id"`
}

func Signup(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var user models.User

	if c.BindJSON(&user) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid data", "data": user})
		c.Abort()
		return
	}

	if !db.NewRecord(user) {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "User could not be created"})
		c.Abort()
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword[:])

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already used"})
		c.Abort()
		return
	}

	// create a customer on kong
	url := "http://apigw-admin.pong.com/consumers"
	customer := Customer{fmt.Sprint(user.ID)}
	pbytes, _ := json.Marshal(customer)
	buff := bytes.NewBuffer(pbytes)

	if _, err := http.Post(url, "application/json", buff); err != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusBadGateway, gin.H{"message": "Something wrong"})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func Login(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)

	var user models.User

	if c.BindJSON(&user) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid data", "data": user})
		c.Abort()
		return
	}

	var matchedUser models.User
	db.Where("email = ?", user.Email).First(&matchedUser)

	if err := bcrypt.CompareHashAndPassword([]byte(matchedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password mismatch"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user logged in"})
}

func UserInfo(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	id := c.Param("id")

	var user models.User
	db.Where("id = ?", id).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
