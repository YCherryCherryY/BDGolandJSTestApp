package main

import (
	"backend/database"
	"backend/handlers"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&models.ContainerStatus{})
	router := gin.Default()
	router.GET("/statuses", handlers.GetStatuses)
	router.POST("/statuses", handlers.AddStatus)

	router.Run(":8080")
}
