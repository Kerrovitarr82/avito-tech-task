package services

import (
	"avito-tech-task/internal/dto"
	"avito-tech-task/internal/models"
	"errors"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) SetIsActive(userID string, isActive bool) (*models.User, error) {
	var user models.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("NOT_FOUND")
		}
		return nil, err
	}

	user.IsActive = isActive
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Team").Where("id = ?", userID).First(&user)
	return &user, nil
}

func (s *UserService) GetUserReviews(userID string) ([]models.PR, error) {
	var user models.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Если пользователь не найден, возвращаем пустой список
			return []models.PR{}, nil
		}
		return nil, err
	}

	var prs []models.PR
	err = s.db.Joins("JOIN reviewers ON reviewers.pr_id = prs.id").
		Where("reviewers.user_id = ?", userID).
		Find(&prs).Error
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (s *UserService) GetUserPRStats(userID string) (*dto.UserPRStats, error) {
	var user models.User
	err := s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("NOT_FOUND")
		}
		return nil, err
	}

	type Result struct {
		Status string
		Count  int64
	}

	var results []Result
	err = s.db.Model(&models.PR{}).
		Select("prs.status, COUNT(DISTINCT prs.id) as count").
		Joins("JOIN reviewers ON reviewers.pr_id = prs.id").
		Where("reviewers.user_id = ?", userID).
		Group("prs.status").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := &dto.UserPRStats{
		UserID:   user.ID,
		Username: user.Username,
	}

	for _, r := range results {
		stats.TotalPRs += r.Count
		if r.Status == "OPEN" {
			stats.OpenPRs = r.Count
		} else if r.Status == "MERGED" {
			stats.MergedPRs = r.Count
		}
	}

	return stats, nil
}
