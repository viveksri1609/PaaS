package models

import "time"

type App struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Status      string    `json:"status"`
	Health      string    `json:"health"`
	ContainerID string    `json:"container_id"`
	Port        int       `json:"port"`
	CreatedAt   time.Time `json:"created_at"`
}
