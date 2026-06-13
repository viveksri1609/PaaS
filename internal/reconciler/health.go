package reconciler

import (
	"fmt"

	"PaaS/internal/db"
	dockerRuntime "PaaS/internal/docker"
	"PaaS/internal/models"
)

func CheckHealth() {

	var apps []models.App

	db.DB.
		Where("status = ?", "running").
		Find(&apps)

	for _, app := range apps {

		if app.ContainerID == "" {
			app.Health = "unhealthy"
			db.DB.Save(&app)
			continue
		}

		running, err := dockerRuntime.IsContainerRunning(
			app.ContainerID,
		)

		if err != nil {

			fmt.Println(
				"health check failed:",
				app.Name,
			)

			app.Health = "unhealthy"

			db.DB.Save(&app)

			continue
		}

		if running {

			app.Health = "healthy"

		} else {

			app.Health = "unhealthy"
		}

		db.DB.Save(&app)
	}
}
