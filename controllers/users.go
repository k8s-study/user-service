package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	url := fmt.Sprintf("%s/consumers", os.Getenv("KONG_HOST"))
	consumer := Consumer{fmt.Sprint(user.ID)}
	pbytes, _ := json.Marshal(consumer)
	buff := bytes.NewBuffer(pbytes)

	resp1, err1 := http.Post(url, "application/json", buff)
	if err1 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusBadGateway, gin.H{"message": "Something wrong"})
		c.Abort()
		return
	}

	defer resp1.Body.Close()

	var data1 RichConsumer
	err1 = json.NewDecoder(resp1.Body).Decode(&data1)
	if err1 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusBadGateway, gin.H{"message": "Something wrong"})
		c.Abort()
		return
	}

	// issue auth token
	tokenUrl := fmt.Sprintf("%s/consumers/%s/key-auth", os.Getenv("KONG_HOST"), data1.Id)
	consumer2 := Empty{}
	pbytes2, _ := json.Marshal(consumer2)
	buff2 := bytes.NewBuffer(pbytes2)
	resp2, err2 := http.Post(tokenUrl, "application/json", buff2)
	if err2 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusBadGateway, gin.H{"message": "Something wrong"})
		c.Abort()
		return
	}

	defer resp2.Body.Close()

	var data2 RichConsumer
	err2 = json.NewDecoder(resp2.Body).Decode(&data2)
	if err2 != nil {
		// FIXME: rollback user from database
		c.JSON(http.StatusBadGateway, gin.H{"message": "Something wrong"})
		c.Abort()
		return
	}

	db.Model(&user).Update("kong_id", data1.Id)

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
