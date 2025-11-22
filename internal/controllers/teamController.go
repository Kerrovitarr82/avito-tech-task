package controllers

import (
	"avito-tech-task/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TeamController struct {
	teamService *services.TeamService
}

func NewTeamController(teamService *services.TeamService) *TeamController {
	return &TeamController{
		teamService: teamService,
	}
}

// AddTeam создает команду с участниками
// POST /team/add
func (c *TeamController) AddTeam(ctx *gin.Context) {
	var req struct {
		TeamName string           `json:"team_name" binding:"required"`
		Members  []map[string]any `json:"members" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	team, err := c.teamService.CreateTeam(req.TeamName, req.Members)
	if err != nil {
		if err.Error() == "TEAM_EXISTS" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	members := make([]gin.H, 0)
	for _, user := range team.Users {
		members = append(members, gin.H{
			"user_id":   user.ID,
			"username":  user.Username,
			"is_active": user.IsActive,
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"team": gin.H{
			"team_name": team.Name,
			"members":   members,
		},
	})
}

// GetTeam возвращает команду с участниками
// GET /team/get?team_name=<name>
func (c *TeamController) GetTeam(ctx *gin.Context) {
	teamName := ctx.Query("team_name")
	if teamName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "team_name is required",
			},
		})
		return
	}

	team, err := c.teamService.GetTeam(teamName)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "team not found",
				},
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	members := make([]gin.H, 0)
	for _, user := range team.Users {
		members = append(members, gin.H{
			"user_id":   user.ID,
			"username":  user.Username,
			"is_active": user.IsActive,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"team_name": team.Name,
		"members":   members,
	})
}

// DeactivateMembers деактивирует указанных участников команды и переназначает их PR
// POST /team/deactivateMembers
func (c *TeamController) DeactivateMembers(ctx *gin.Context) {
	var req struct {
		TeamName string   `json:"team_name" binding:"required"`
		UserIDs  []string `json:"user_ids" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	deactivatedCount, reassignedCount, err := c.teamService.DeactivateTeamMembers(req.TeamName, req.UserIDs)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "team not found",
				},
			})
		case "NO_VALID_USERS":
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "NO_VALID_USERS",
					"message": "no valid users found in the specified team",
				},
			})
		case "INVALID_REQUEST":
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "user_ids cannot be empty",
				},
			})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"team_name":          req.TeamName,
		"deactivated_users":  deactivatedCount,
		"reassigned_reviews": reassignedCount,
	})
}
