package main

import (
	"github.com/gin-gonic/gin"
	"github.com/k8s-study/user-service/db"
	"github.com/k8s-study/user-service/controllers"
)

func main() {
	db.Init()

	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/health", controllers.Health)
	}
	r.Run()
}
