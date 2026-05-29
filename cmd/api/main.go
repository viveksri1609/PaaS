package main

import (
	"github.com/gin-gonic/gin"

	"mini-paas/internal/db"
	"mini-paas/internal/handlers"
)

func main() {
	db.Connect()

	router := gin.Default()

	router.POST("/apps", handlers.CreateApp)
	router.GET("/apps", handlers.GetApps)

	router.Run(":8081")
}
