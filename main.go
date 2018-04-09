package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k8s-study/user-service/db"
	"github.com/k8s-study/user-service/controllers"
)

func main() {
	db.Init()

	r := gin.Default()

	r.Use(db.Init())

	v1 := r.Group("/v1")
	{
		v1.GET("/health", controllers.Health)
		v1.POST("/signup", controllers.Signup)
		v1.POST("/login", controllers.Login)
		v1.GET("/users/:id", controllers.UserInfo)
	}
	r.Run()
}
