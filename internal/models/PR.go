package models

import (
	"gorm.io/gorm"
	"time"
)

type PR struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	AuthorID  uint           `gorm:"not null;index" json:"authorId"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	Status    string         `gorm:"not null" json:"status"` // ограничения прописаны в коде
	Reviewers []User         `gorm:"many2many:reviewers" json:"reviewers"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
