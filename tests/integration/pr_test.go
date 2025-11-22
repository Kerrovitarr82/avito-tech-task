package integration

import (
	"avito-tech-task/internal/models"
	"avito-tech-task/internal/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPRService_CreatePR_Integration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)
	SeedTestData(t, db)

	prService := services.NewPRService(db)

	t.Run("should create PR successfully with reviewers", func(t *testing.T) {
		pr, err := prService.CreatePR("pr-1", "Feature: Add login", "test-user-1")

		assert.NoError(t, err)
		assert.NotNil(t, pr)
		assert.Equal(t, "pr-1", pr.ID)
		assert.Equal(t, "Feature: Add login", pr.Name)
		assert.Equal(t, "test-user-1", pr.AuthorID)
		assert.Equal(t, "OPEN", pr.Status)
		assert.NotZero(t, pr.CreatedAt)
		assert.LessOrEqual(t, len(pr.Reviewers), 2)
		assert.GreaterOrEqual(t, len(pr.Reviewers), 0)

		for _, reviewer := range pr.Reviewers {
			assert.NotEqual(t, "test-user-1", reviewer.ID)
		}

		var dbPR models.PR
		db.Preload("Author").Preload("Reviewers").First(&dbPR, "id = ?", "pr-1")
		assert.Equal(t, "pr-1", dbPR.ID)
		assert.Equal(t, "test-user-1", dbPR.Author.ID)
	})
}
