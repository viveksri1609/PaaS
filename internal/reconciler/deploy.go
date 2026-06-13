package reconciler

import (
	"fmt"

	"PaaS/internal/db"
	dockerRuntime "PaaS/internal/docker"
	"PaaS/internal/models"
)

func DeployPendingApps() {

	var apps []models.App

	db.DB.
		Where("status = ?", "pending").
		Find(&apps)

	for _, app := range apps {

		fmt.Println(
			"deploying:",
			app.Name,
		)

		containerID, err :=
			dockerRuntime.RunContainer(
				app.Name,
				app.Name,
				app.Image,
			)

		if err != nil {

			app.Status = "failed"
			app.Health = "unknown"

		} else {

			app.Status = "running"
			app.Health = "healthy"
			app.ContainerID = containerID
		}

		db.DB.Save(&app)
	}
}
