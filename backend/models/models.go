package models

import (
	"time"

	"gorm.io/gorm"
)

type ContainerStatus struct {
	gorm.Model
	IP          string    `json:"ip"`
	Status      string    `json:"status"`
	Time        time.Time `json:"time"`
	TimeSuccses time.Time `json:"succsesTime"`
}
