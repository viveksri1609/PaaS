package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"PaaS/internal/db"
	"PaaS/internal/models"
)

type ScaleRequest struct {
	Replicas int `json:"replicas"`
}

func ScaleApp(c *gin.Context) {

	id := c.Param("id")

	var app models.App

	if err := db.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "app not found",
		})
		return
	}

	var req ScaleRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Replicas < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "replicas must be greater than or equal to zero",
		})
		return
	}

	app.Replicas = req.Replicas

	if err := db.DB.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, app)
}
