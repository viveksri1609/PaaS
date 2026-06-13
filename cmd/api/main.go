package main

import (
	"github.com/gin-gonic/gin"

	"PaaS/internal/db"
	"PaaS/internal/handlers"
)

func main() {
	db.Connect()

	router := gin.Default()

	router.POST("/apps", handlers.CreateApp)
	router.GET("/apps", handlers.GetApps)
	router.DELETE("/apps/:id", handlers.DeleteApp)
	router.POST(
		"/apps/:id/scale",
		handlers.ScaleApp,
	)
	router.Run(":8081")
}
