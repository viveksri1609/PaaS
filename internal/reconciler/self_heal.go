package reconciler

import (
	"fmt"

	"PaaS/internal/db"
	dockerRuntime "PaaS/internal/docker"
	"PaaS/internal/models"
)

func HealUnhealthyApps() {

	var apps []models.App

	db.DB.
		Where("health = ?", "unhealthy").
		Find(&apps)

	for _, app := range apps {

		fmt.Println(
			"recovering app:",
			app.Name,
		)

		if app.ContainerID != "" {
			if err := dockerRuntime.DeleteContainer(app.ContainerID); err != nil {
				fmt.Println(
					"failed to remove unhealthy container:",
					app.Name,
					err,
				)
			}
		}

		containerID, err :=
			dockerRuntime.RunContainer(
				app.Name,
				app.Image,
			)

		if err != nil {

			fmt.Println(
				"failed to recover:",
				app.Name,
				err,
			)

			continue
		}

		app.ContainerID = containerID
		app.Health = "healthy"

		db.DB.Save(&app)

		fmt.Println(
			"app recovered:",
			app.Name,
		)
	}
}
