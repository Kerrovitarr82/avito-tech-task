package database

import (
	"avito-tech-task/internal/models"
	"avito-tech-task/internal/services"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func SeedDB(db *gorm.DB, prService *services.PRService) {
	var count int64
	db.Model(&models.Team{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return
	}

	teams := make([]models.Team, 20)
	for i := 0; i < 20; i++ {
		teams[i] = models.Team{
			Name: fmt.Sprintf("Team-%d", i+1),
		}
	}

	if err := db.Create(&teams).Error; err != nil {
		log.Fatalf("Failed to seed teams: %v", err)
	}

	users := make([]models.User, 200)
	for i := 0; i < 200; i++ {
		teamIndex := i / 10
		users[i] = models.User{
			ID:       fmt.Sprintf("u%d", i+1),
			Username: fmt.Sprintf("user%d", i+1),
			IsActive: true,
			TeamName: teams[teamIndex].Name,
		}
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	for i := 0; i < 100; i++ {
		_, err := prService.CreatePR(fmt.Sprintf("pr-%d", i+1), fmt.Sprintf("pr_title-%d", i+1), fmt.Sprintf("u%d", i+2))
		if err != nil {
			log.Fatalf("Failed to seed pull requests: %v", err)
		}
	}

	log.Println("Database seeded successfully!")
}
