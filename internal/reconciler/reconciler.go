package reconciler

import "time"

func Start() {

	for {

		DeployPendingApps()

		CheckHealth()

		HealUnhealthyApps()

		time.Sleep(
			5 * time.Second,
		)
	}
}
