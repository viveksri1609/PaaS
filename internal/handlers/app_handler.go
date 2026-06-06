package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"PaaS/internal/db"
	dockerRuntime "PaaS/internal/docker"
	"PaaS/internal/models"
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

func DeleteApp(c *gin.Context) {

	id := c.Param("id")

	var app models.App

	err := db.DB.First(&app, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "app not found",
		})
		return
	}

	if app.ContainerID != "" {

		err := dockerRuntime.DeleteContainer(app.ContainerID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	db.DB.Delete(&app)

	c.JSON(http.StatusOK, gin.H{
		"message": "app deleted successfully",
	})
}
