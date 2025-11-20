package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name" validate:"required,min=2,max=100"`
	IsActive  bool           `gorm:"not null;default:true" json:"isActive"`
	TeamName  string         `gorm:"not null;index" json:"teamName"`
	Team      Team           `gorm:"foreignKey:TeamName;references:Name" json:"team"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
