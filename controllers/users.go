package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/k8s-study/user-service/client"
	"github.com/k8s-study/user-service/models"
	"golang.org/x/crypto/bcrypt"
)

type Consumer struct {
	Custom_id string `json:"custom_id"`
}

type RichConsumer struct {
	Id       string `json:"id"`
	CustomId string `json:"custom_id"`
	Key      string `json:"key"`
}

type Empty struct {
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

	// create a consumer on kong
	consumer1 := Consumer{fmt.Sprint(user.ID)}
	client := client.NewClient(c.Request)
	req1, err1 := client.NewRequest("POST", "/consumers", consumer1)
	if err1 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusInternalServerError, gin.H{"message": err1.Error()})
		c.Abort()
		return
	}

	var data1 RichConsumer
	resp1, err1 := client.Do(req1, &data1)
	if err1 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Fail to create a new consumer on kong"})
		c.Abort()
		return
	}
	if resp1.StatusCode >= 400 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": data1.CustomId})
		c.Abort()
		return
	}

	// issue auth token
	consumer2 := Empty{}
	req2, err2 := client.NewRequest("POST", fmt.Sprintf("/consumers/%s/key-auth", data1.Id), consumer2)
	if err2 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		c.Abort()
		return
	}

	var data2 RichConsumer
	_, err2 = client.Do(req2, &data2)
	if err2 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Fail to get auth token from kong"})
		c.Abort()
		return
	}

	db.Model(&user).Update("kong_id", data1.Id)
	if result := db.Model(&user).Update("kong_id", data1.Id); result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Fail to update kong_id"})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"key":     data2.Key,
	})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Password mismatch"})
		c.Abort()
		return
	}

	// issue auth token
	client := client.NewClient(c.Request)
	req, err := client.NewRequest("POST", fmt.Sprintf("/consumers/%s/key-auth", matchedUser.KongId), Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	var data RichConsumer
	_, err = client.Do(req, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Fail to get auth token from kong"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user logged in",
		"key":     data.Key,
	})
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

func CurrentUserInfo(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	id := c.Request.Header.Get("X-Consumer-Custom-ID")

	var user models.User
	db.Where("id = ?", id).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
