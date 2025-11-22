package integration

import (
	"avito-tech-task/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		"db_test", "postgres", "test_pass", "test_db", "disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Team{}, &models.User{}, &models.PR{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
	db.Exec("TRUNCATE TABLE reviewers CASCADE")
	db.Exec("TRUNCATE TABLE prs CASCADE")
	db.Exec("TRUNCATE TABLE users CASCADE")
	db.Exec("TRUNCATE TABLE teams CASCADE")
}

func SeedTestData(t *testing.T, db *gorm.DB) {
	team := models.Team{Name: "TestTeam"}
	db.Create(&team)

	users := []models.User{
		{ID: "test-user-1", Username: "testuser1", IsActive: true, TeamName: "TestTeam"},
		{ID: "test-user-2", Username: "testuser2", IsActive: true, TeamName: "TestTeam"},
		{ID: "test-user-3", Username: "testuser3", IsActive: false, TeamName: "TestTeam"},
	}
	db.Create(&users)
}
