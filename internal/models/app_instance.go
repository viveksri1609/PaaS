package models

import "time"

type AppInstance struct {
	ID          uint `gorm:"primaryKey"`
	AppID       uint
	ContainerID string
	Status      string
	CreatedAt   time.Time
}
