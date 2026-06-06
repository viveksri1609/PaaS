package main

import (
	"PaaS/internal/db"
	"PaaS/internal/reconciler"
)

func main() {
	db.Connect()

	reconciler.Start()
}
