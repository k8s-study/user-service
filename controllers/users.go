package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/k8s-study/user-service/models"
	"net/http"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

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
