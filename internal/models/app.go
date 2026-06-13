package models

import "time"

type App struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Image       string
	Status      string
	Health      string
	ContainerID string

	Replicas int

	CreatedAt time.Time
}
