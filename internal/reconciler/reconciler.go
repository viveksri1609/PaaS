package reconciler

import (
	"fmt"
	"time"

	"mini-paas/internal/db"
	dockerRuntime "mini-paas/internal/docker"
	"mini-paas/internal/models"
)

func Start() {
	for {
		var apps []models.App

		db.DB.Where("status = ?", "pending").Find(&apps)

		for _, app := range apps {
			fmt.Println("deploying app", app.Name)

			containerID, err := dockerRuntime.RunContainer(app.Name, app.Image)

			if err != nil {
				fmt.Println(err)
				app.Status = "failed"
			} else {
				app.Status = "running"
				app.ContainerID = containerID
			}

			db.DB.Save(&app)
		}

		time.Sleep(5 * time.Second)
	}
}
