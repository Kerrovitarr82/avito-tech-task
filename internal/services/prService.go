package services

import (
	"avito-tech-task/internal/models"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type PRService struct {
	db *gorm.DB
}

func NewPRService(db *gorm.DB) *PRService {
	return &PRService{db: db}
}

func (s *PRService) CreatePR(prID string, title string, authorID string) (*models.PR, error) {
	var existingPR models.PR
	err := s.db.Where("id = ?", prID).First(&existingPR).Error
	if err == nil {
		return nil, errors.New("PR_EXISTS")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var author models.User
	err = s.db.Preload("Team").Where("id = ?", authorID).First(&author).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("NOT_FOUND")
		}
		return nil, err
	}

	pr := models.PR{
		ID:       prID,
		Name:     title,
		AuthorID: authorID,
		Status:   "OPEN",
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&pr).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	reviewers, err := s.selectReviewers(tx, authorID, author.TeamName, 2)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(reviewers) > 0 {
		if err := tx.Model(&pr).Association("Reviewers").Append(reviewers); err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	s.db.Preload("Author").Preload("Reviewers").Where("id = ?", pr.ID).First(&pr)
	return &pr, nil
}

func (s *PRService) selectReviewers(db *gorm.DB, authorID string, teamName string, maxCandidates int) ([]models.User, error) {
	var candidates []models.User
	err := db.Where("team_name = ? AND id != ? AND is_active = ?", teamName, authorID, true).
		Find(&candidates).Error
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	if len(candidates) < maxCandidates {
		maxCandidates = len(candidates)
	}
	return candidates[:maxCandidates], nil
}

func (s *PRService) MergePR(prID string) (*models.PR, error) {
	var pr models.PR
	err := s.db.Preload("Author").Preload("Reviewers").Where("id = ?", prID).First(&pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("NOT_FOUND")
		}
		return nil, err
	}

	if pr.Status == "MERGED" {
		return &pr, nil
	}

	pr.Status = "MERGED"
	pr.MergedAt = time.Now()
	if err := s.db.Save(&pr).Error; err != nil {
		return nil, err
	}

	return &pr, nil
}

func (s *PRService) ReassignReviewer(prID string, oldUserID string) (*models.PR, string, error) {
	var pr models.PR
	var oldUser *models.User
	err := s.db.Where("id = ?", prID).First(&pr).Error
	err = s.db.Preload("Author").Preload("Author.Team").Preload("Reviewers").Where("id = ?", prID).First(&pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("NOT_FOUND")
		}
		return nil, "", err
	}

	if pr.Status == "MERGED" {
		return nil, "", errors.New("PR_MERGED")
	}

	isAssigned := false
	for i := range pr.Reviewers {
		if pr.Reviewers[i].ID == oldUserID {
			oldUser = &pr.Reviewers[i]
			isAssigned = true
			break
		}
	}

	if !isAssigned {
		return nil, "", errors.New("NOT_ASSIGNED")
	}

	reviewers, err := s.selectReviewers(s.db, oldUserID, oldUser.TeamName, 1)
	if err != nil {
		return nil, "", errors.New("NO_CANDIDATE")
	}

	newReviewer := reviewers[0]

	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, "", tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&pr).Association("Reviewers").Delete(oldUser); err != nil {
		tx.Rollback()
		return nil, "", err
	}

	if err := tx.Model(&pr).Association("Reviewers").Append(&newReviewer); err != nil {
		tx.Rollback()
		return nil, "", err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", err
	}

	s.db.Preload("Author").Preload("Reviewers").Where("id = ?", pr.ID).First(&pr)

	return &pr, newReviewer.ID, nil
}
