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
		Name:     req.Name,
		Image:    req.Image,
		Status:   "pending",
		Health:   "unknown",
		Replicas: 1,
	}

	if err := db.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, app)
}

func GetApps(c *gin.Context) {
	var apps []models.App

	db.DB.Find(&apps)

	c.JSON(http.StatusOK, apps)
}

type ContainerLog struct {
	ContainerID string `json:"container_id"`
	Logs        string `json:"logs,omitempty"`
	Error       string `json:"error,omitempty"`
}

type ContainerMetric struct {
	ContainerID   string  `json:"container_id"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   uint64  `json:"memory_usage"`
	MemoryLimit   uint64  `json:"memory_limit"`
	MemoryPercent float64 `json:"memory_percent"`
	Restarts      int64   `json:"restarts"`
	Status        string  `json:"status"`
	Error         string  `json:"error,omitempty"`
}

func GetAppLogs(c *gin.Context) {
	id := c.Param("id")

	var app models.App
	if err := db.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	containers := []string{}
	if app.ContainerID != "" {
		containers = append(containers, app.ContainerID)
	}

	var instances []models.AppInstance
	db.DB.Where("app_id = ?", app.ID).Find(&instances)
	for _, instance := range instances {
		containers = append(containers, instance.ContainerID)
	}

	if len(containers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no containers found for app"})
		return
	}

	tail := c.DefaultQuery("tail", "100")
	logs := make([]ContainerLog, 0, len(containers))

	for _, containerID := range containers {
		body, err := dockerRuntime.ContainerLogs(containerID, tail)
		entry := ContainerLog{ContainerID: containerID}
		if err != nil {
			entry.Error = err.Error()
		} else {
			entry.Logs = string(body)
		}
		logs = append(logs, entry)
	}

	c.JSON(http.StatusOK, gin.H{
		"app_id": app.ID,
		"logs":   logs,
	})
}

func GetAppMetrics(c *gin.Context) {
	id := c.Param("id")

	var app models.App
	if err := db.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
		return
	}

	containers := []string{}
	if app.ContainerID != "" {
		containers = append(containers, app.ContainerID)
	}

	var instances []models.AppInstance
	db.DB.Where("app_id = ?", app.ID).Find(&instances)
	for _, instance := range instances {
		containers = append(containers, instance.ContainerID)
	}

	if len(containers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no containers found for app"})
		return
	}

	metrics := make([]ContainerMetric, 0, len(containers))

	for _, containerID := range containers {
		stats, statsErr := dockerRuntime.ContainerStats(containerID)
		inspect, inspectErr := dockerRuntime.ContainerInspect(containerID)

		entry := ContainerMetric{ContainerID: containerID}
		if inspectErr != nil {
			entry.Error = inspectErr.Error()
			metrics = append(metrics, entry)
			continue
		}

		if inspect.State != nil {
			entry.Status = string(inspect.State.Status)
			entry.Restarts = int64(inspect.RestartCount)
		}

		if statsErr != nil {
			metrics = append(metrics, entry)
			continue
		}

		cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
		if systemDelta > 0 && cpuDelta > 0 {
			cpuCount := len(stats.CPUStats.CPUUsage.PercpuUsage)
			if cpuCount == 0 {
				cpuCount = 1
			}
			entry.CPUPercent = cpuDelta / systemDelta * float64(cpuCount) * 100.0
		}

		memoryUsage := stats.MemoryStats.Usage
		if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
			if memoryUsage > cache {
				memoryUsage -= cache
			}
		}
		entry.MemoryUsage = memoryUsage
		entry.MemoryLimit = stats.MemoryStats.Limit
		if stats.MemoryStats.Limit > 0 {
			entry.MemoryPercent = float64(memoryUsage) / float64(stats.MemoryStats.Limit) * 100.0
		}

		metrics = append(metrics, entry)
	}

	c.JSON(http.StatusOK, gin.H{
		"app_id":  app.ID,
		"metrics": metrics,
	})
}

func DeleteApp(c *gin.Context) {

	id := c.Param("id")

	var app models.App

	if err := db.DB.First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "app not found",
		})
		return
	}

	var instances []models.AppInstance

	db.DB.
		Where("app_id = ?", app.ID).
		Find(&instances)

	for _, instance := range instances {

		err := dockerRuntime.DeleteContainer(
			instance.ContainerID,
		)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		db.DB.Delete(&instance)
	}

	if app.ContainerID != "" {
		if err := dockerRuntime.DeleteContainer(app.ContainerID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	if err := db.DB.Delete(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "app deleted successfully",
	})
}
