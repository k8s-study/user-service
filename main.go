package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k8s-study/user-service/controllers"
	"github.com/k8s-study/user-service/db"
)

func main() {
	db.Init()

	r := gin.Default()

	r.Use(db.Init())

	r.GET("/health", controllers.Health)

	inV1 := r.Group("/private/v1")
	{
		inV1.GET("/users/:id", controllers.UserInfo)
	}

	exV1 := r.Group("/public/v1")
	{
		exV1.GET("/user", controllers.CurrentUserInfo)
		exV1.GET("/users/:id", controllers.UserInfo)
		exV1.POST("/signup", controllers.Signup)
		exV1.POST("/login", controllers.Login)
	}

	r.Run()
}
