package reconciler

import "time"

func Start() {

	for {

		DeployPendingApps()

		CheckHealth()

		HealUnhealthyApps()

		ReconcileReplicas()

		time.Sleep(
			5 * time.Second,
		)
	}
}
