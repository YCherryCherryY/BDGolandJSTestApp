package handlers

import (
	"backend/database"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetStatuses(context *gin.Context) {
	var statuses []models.ContainerStatus
	database.DB.Find(&statuses)
	context.JSON(http.StatusOK, statuses)
}

func AddStatus(context *gin.Context) {
	var status models.ContainerStatus
	if err := context.ShouldBindJSON(&status); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	checkTime, err := time.Parse(time.RFC3339, context.PostForm("check_time"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid check_time"})
		return
	}
	status.Time = checkTime

	lastSuccessTime, err := time.Parse(time.RFC3339, context.PostForm("last_success_time"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid last_success_time"})
		return
	}
	status.TimeSuccses = lastSuccessTime

	database.DB.Create(&status)
	context.JSON(http.StatusCreated, status)
}
