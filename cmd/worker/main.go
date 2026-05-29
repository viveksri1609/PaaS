package main

import (
	"mini-paas/internal/db"
	"mini-paas/internal/reconciler"
)

func main() {
	db.Connect()

	reconciler.Start()
}
