package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"PaaS/internal/models"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=admin password=admin dbname=paas port=5432 sslmode=disable"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}

	fmt.Println("database connected")

	database.AutoMigrate(&models.App{})

	DB = database
}
