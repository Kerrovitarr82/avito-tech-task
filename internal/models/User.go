package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"user_id"`
	Username  string         `gorm:"not null" json:"username" validate:"required,min=2,max=100"`
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	TeamName  string         `gorm:"not null;index" json:"team_name"`
	Team      Team           `gorm:"foreignKey:TeamName;references:Name" json:"team"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
