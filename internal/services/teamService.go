package services

import (
	"avito-tech-task/internal/models"
	"errors"
	"gorm.io/gorm"
)

type TeamService struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{db: db}
}

func (s *TeamService) CreateTeam(teamName string, members []map[string]any) (*models.Team, error) {
	var existingTeam models.Team
	err := s.db.Where("name = ?", teamName).First(&existingTeam).Error
	if err == nil {
		return nil, errors.New("TEAM_EXISTS")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if teamName == "" {
		return nil, errors.New("team name is required")
	}
	if len(members) == 0 {
		return nil, errors.New("at least one member is required")
	}

	team := models.Team{
		Name: teamName,
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

	if err := tx.Create(&team).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	users := make([]models.User, 0)
	for _, member := range members {
		userID, ok := member["user_id"].(string)
		if !ok {
			tx.Rollback()
			return nil, errors.New("invalid user_id format")
		}

		username, ok := member["username"].(string)
		if !ok {
			tx.Rollback()
			return nil, errors.New("invalid username format")
		}

		isActive, ok := member["is_active"].(bool)
		if !ok {
			tx.Rollback()
			return nil, errors.New("invalid is_active format")
		}

		if len(username) < 2 || len(username) > 100 {
			tx.Rollback()
			return nil, errors.New("username must be between 2 and 100 characters")
		}

		user := models.User{
			ID:       userID,
			Username: username,
			IsActive: isActive,
			TeamName: teamName,
		}

		var existingUser models.User
		err := tx.Where("id = ?", userID).First(&existingUser).Error
		if err == nil { // обновление существующего юзера
			existingUser.TeamName = teamName
			existingUser.IsActive = isActive
			if err := tx.Save(&existingUser).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			users = append(users, existingUser)
		} else if errors.Is(err, gorm.ErrRecordNotFound) { // создание нового юзера
			if err := tx.Create(&user).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			users = append(users, user)
		} else {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	team.Users = users
	return &team, nil
}

func (s *TeamService) GetTeam(teamName string) (*models.Team, error) {
	var team models.Team
	err := s.db.Preload("Users").Where("name = ?", teamName).First(&team).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("NOT_FOUND")
		}
		return nil, err
	}
	return &team, nil
}

func (s *TeamService) DeactivateTeamMembers(teamName string, userIDs []string) (int, int, error) {
	var team models.Team
	err := s.db.Where("name = ?", teamName).First(&team).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, errors.New("NOT_FOUND")
		}
		return 0, 0, err
	}

	if len(userIDs) == 0 {
		return 0, 0, errors.New("INVALID_REQUEST")
	}

	var validUserIDs []string
	err = s.db.Model(&models.User{}).
		Where("id IN ? AND team_name = ?", userIDs, teamName).
		Pluck("id", &validUserIDs).Error
	if err != nil {
		return 0, 0, err
	}

	if len(validUserIDs) == 0 {
		return 0, 0, errors.New("NO_VALID_USERS")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return 0, 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	reassignedCount := 0

	var prs []models.PR
	err = tx.Preload("Author").
		Preload("Reviewers", "id IN ?", validUserIDs).
		Joins("JOIN reviewers ON reviewers.pr_id = prs.id").
		Where("reviewers.user_id IN ? AND prs.status = ?", validUserIDs, "OPEN").
		Distinct().
		Find(&prs).Error
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	prsByTeam := make(map[string][]models.PR)
	for _, pr := range prs {
		prsByTeam[pr.Author.TeamName] = append(prsByTeam[pr.Author.TeamName], pr)
	}

	for teamName, teamPRs := range prsByTeam {
		var candidates []models.User
		err = tx.Where("team_name = ? AND is_active = ? AND id NOT IN ?",
			teamName, true, validUserIDs).
			Find(&candidates).Error
		if err != nil {
			tx.Rollback()
			return 0, 0, err
		}

		if len(candidates) == 0 {
			continue
		}

		candidateIndex := 0

		for _, pr := range teamPRs {
			for _, reviewer := range pr.Reviewers {
				var selectedCandidate *models.User
				for i := 0; i < len(candidates); i++ {
					candidate := &candidates[(candidateIndex+i)%len(candidates)]
					if candidate.ID != pr.AuthorID {
						selectedCandidate = candidate
						candidateIndex = (candidateIndex + i + 1) % len(candidates)
						break
					}
				}

				if selectedCandidate == nil {
					continue
				}

				if err := tx.Model(&pr).Association("Reviewers").Delete(&reviewer); err != nil {
					tx.Rollback()
					return 0, 0, err
				}

				if err := tx.Model(&pr).Association("Reviewers").Append(selectedCandidate); err != nil {
					tx.Rollback()
					return 0, 0, err
				}

				reassignedCount++
			}
		}
	}

	result := tx.Model(&models.User{}).
		Where("id IN ? AND team_name = ?", validUserIDs, teamName).
		Update("is_active", false)
	if result.Error != nil {
		tx.Rollback()
		return 0, 0, result.Error
	}

	deactivatedCount := int(result.RowsAffected)

	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	return deactivatedCount, reassignedCount, nil
}
