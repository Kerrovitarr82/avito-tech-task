package models

import (
	"gorm.io/gorm"
	"time"
)

type PR struct {
	ID        string         `gorm:"primaryKey" json:"pull_request_id"`
	Name      string         `gorm:"not null" json:"pull_request_name"`
	AuthorID  string         `gorm:"not null;index" json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	Status    string         `gorm:"not null;check:status IN ('OPEN', 'MERGED')" json:"status"`
	Reviewers []User         `gorm:"many2many:reviewers" json:"assigned_reviewers"`
	MergedAt  time.Time      `json:"mergedAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
