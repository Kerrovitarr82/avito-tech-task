package models

import (
	"gorm.io/gorm"
	"time"
)

type Team struct {
	Name      string         `gorm:"primaryKey" json:"team_name"`
	Users     []User         `gorm:"foreignKey:TeamName;references:Name" json:"members"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
