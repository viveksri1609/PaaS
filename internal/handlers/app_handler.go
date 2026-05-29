package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mini-paas/internal/db"
	"mini-paas/internal/models"
)

type CreateAppRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func CreateApp(c *gin.Context) {
	var req CreateAppRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	app := models.App{
		Name:   req.Name,
		Image:  req.Image,
		Status: "pending",
	}

	db.DB.Create(&app)

	c.JSON(http.StatusCreated, app)
}

func GetApps(c *gin.Context) {
	var apps []models.App

	db.DB.Find(&apps)

	c.JSON(http.StatusOK, apps)
}
