package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"PaaS/internal/db"
	"PaaS/internal/handlers"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "check the local API health endpoint")
	flag.Parse()
	if *healthcheck {
		client := http.Client{Timeout: 2 * time.Second}
		response, err := client.Get("http://127.0.0.1:8081/healthz")
		if err != nil {
			os.Exit(1)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		return
	}

	db.Connect()

	router := gin.Default()

	router.GET("/healthz", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})
	router.POST("/apps", handlers.CreateApp)
	router.GET("/apps", handlers.GetApps)
	router.GET("/apps/:id/logs", handlers.GetAppLogs)
	router.GET("/apps/:id/metrics", handlers.GetAppMetrics)
	router.DELETE("/apps/:id", handlers.DeleteApp)
	router.POST(
		"/apps/:id/scale",
		handlers.ScaleApp,
	)
	router.Run(":8081")
}
