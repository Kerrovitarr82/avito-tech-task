package integration

import (
	"avito-tech-task/internal/models"
	"avito-tech-task/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_SetIsActive_Integration(t *testing.T) {
	// Setup
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)
	SeedTestData(t, db)

	userService := services.NewUserService(db)

	t.Run("should activate user successfully", func(t *testing.T) {
		// Act
		user, err := userService.SetIsActive("test-user-3", true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.True(t, user.IsActive)

		// Проверяем в БД
		var dbUser models.User
		db.First(&dbUser, "id = ?", "test-user-3")
		assert.True(t, dbUser.IsActive)
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		// Act
		user, err := userService.SetIsActive("non-existent", true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})
}

func TestUserService_GetUserPRStats_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)
	SeedTestData(t, db)

	user := []models.User{{ID: "test-user-1", Username: "testuser1", IsActive: true, TeamName: "TestTeam"}}

	prs := []models.PR{
		{ID: "pr-1", Name: "PR 1", AuthorID: "test-user-2", Status: "OPEN", Reviewers: user},
		{ID: "pr-2", Name: "PR 2", AuthorID: "test-user-2", Status: "OPEN", Reviewers: user},
		{ID: "pr-3", Name: "PR 3", AuthorID: "test-user-2", Status: "MERGED", Reviewers: user},
	}
	db.Create(&prs)

	userService := services.NewUserService(db)

	t.Run("should return correct PR stats", func(t *testing.T) {
		// Act
		stats, err := userService.GetUserPRStats("test-user-1")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, "test-user-1", stats.UserID)
		assert.Equal(t, "testuser1", stats.Username)
		assert.Equal(t, int64(3), stats.TotalPRs)
		assert.Equal(t, int64(2), stats.OpenPRs)
		assert.Equal(t, int64(1), stats.MergedPRs)
	})
}
