package reconciler

import (
	"PaaS/internal/db"
	"PaaS/internal/models"
)

func ReconcileReplicas() {

	var apps []models.App

	db.DB.Find(&apps)

	for _, app := range apps {

		var instances []models.AppInstance

		db.DB.
			Where("app_id = ?", app.ID).
			Find(&instances)

		actual := len(instances)
		if app.ContainerID != "" {
			actual++
		}

		if actual < app.Replicas {

			scaleUp(
				&app,
				app.Replicas-actual,
			)

		} else if actual > app.Replicas {

			scaleDown(
				&app,
				actual-app.Replicas,
			)
		}
	}
}
