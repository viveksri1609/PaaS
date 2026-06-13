package reconciler

import (
	"PaaS/internal/db"
	dockerRuntime "PaaS/internal/docker"
	"PaaS/internal/models"
)

func scaleDown(
	app *models.App,
	count int,
) {
	if count <= 0 {
		return
	}

	var instances []models.AppInstance

	db.DB.
		Where("app_id = ?", app.ID).
		Limit(count).
		Find(&instances)

	for _, instance := range instances {

		_ = dockerRuntime.DeleteContainer(
			instance.ContainerID,
		)

		db.DB.Delete(&instance)
		count--

		if count == 0 {
			return
		}
	}

	if count > 0 && app.ContainerID != "" {
		_ = dockerRuntime.DeleteContainer(app.ContainerID)
		app.ContainerID = ""
		app.Status = "stopped"
		app.Health = "unknown"
		db.DB.Save(app)
	}
}
