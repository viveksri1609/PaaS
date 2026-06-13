package reconciler

import (
	"fmt"
	"time"

	"PaaS/internal/db"
	dockerRuntime "PaaS/internal/docker"
	"PaaS/internal/models"
)

func scaleUp(
	app *models.App,
	count int,
) {
	if count <= 0 {
		return
	}

	if app.ContainerID == "" {
		containerID, err := dockerRuntime.RunContainer(
			app.Name,
			app.Name,
			app.Image,
		)
		if err == nil {
			app.ContainerID = containerID
			app.Status = "running"
			app.Health = "healthy"
			db.DB.Save(app)
			count--
		}
	}

	for i := 0; i < count; i++ {

		containerID, err :=
			dockerRuntime.RunContainer(
				app.Name,
				fmt.Sprintf(
					"%s-%d-%d",
					app.Name,
					time.Now().UnixNano(),
					i,
				),
				app.Image,
			)

		if err != nil {
			continue
		}

		instance := models.AppInstance{
			AppID:       app.ID,
			ContainerID: containerID,
			Status:      "running",
		}

		db.DB.Create(&instance)
	}
}
