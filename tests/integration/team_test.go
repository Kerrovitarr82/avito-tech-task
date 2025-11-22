// tests/integration/team_test.go
package integration

import (
	"avito-tech-task/internal/models"
	"avito-tech-task/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamService_DeactivateTeamMembers_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)
	SeedTestData(t, db)

	// Создаем PR с ревьюверами
	var user1, user2 models.User
	db.First(&user1, "id = ?", "test-user-1")
	db.First(&user2, "id = ?", "test-user-2")

	pr := models.PR{
		ID:       "pr-test",
		Name:     "Test PR",
		AuthorID: "test-user-1",
		Status:   "OPEN",
	}
	db.Create(&pr)
	db.Model(&pr).Association("Reviewers").Append(&user2)

	teamService := services.NewTeamService(db)

	t.Run("should deactivate users and reassign reviewers", func(t *testing.T) {
		// Act
		deactivated, reassigned, err := teamService.DeactivateTeamMembers(
			"TestTeam",
			[]string{"test-user-2"},
		)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, reassigned)
		assert.Equal(t, 1, deactivated)

		// Проверяем деактивацию в БД
		var user models.User
		db.First(&user, "id = ?", "test-user-2")
		assert.False(t, user.IsActive)

		// Проверяем переназначение ревьювера
		var updatedPR models.PR
		db.Preload("Reviewers").First(&updatedPR, "id = ?", "pr-test")

		// user-2 не должен быть в ревьюверах
		for _, reviewer := range updatedPR.Reviewers {
			assert.NotEqual(t, "test-user-2", reviewer.ID)
		}
	})

	t.Run("should return error for non-existent team", func(t *testing.T) {
		// Act
		_, _, err := teamService.DeactivateTeamMembers(
			"NonExistentTeam",
			[]string{"test-user-1"},
		)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "NOT_FOUND", err.Error())
	})
}
