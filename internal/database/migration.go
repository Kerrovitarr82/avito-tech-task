package database

import (
	"avito-tech-task/internal/models"
	"gorm.io/gorm"
	"log"
)

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&models.Team{}, &models.User{}, &models.PR{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration succeeded!")
}
