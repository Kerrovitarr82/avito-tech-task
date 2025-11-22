package controllers

import (
	"avito-tech-task/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// SetIsActive устанавливает флаг активности пользователя
// POST /users/setIsActive
func (c *UserController) SetIsActive(ctx *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		IsActive bool   `json:"is_active"`
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

	user, err := c.userService.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "user not found",
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

	ctx.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"user_id":   user.ID,
			"username":  user.Username,
			"team_name": user.TeamName,
			"is_active": user.IsActive,
		},
	})
}

// GetReview возвращает PR'ы, где пользователь назначен ревьювером
// GET /users/getReview?user_id=<id>
func (c *UserController) GetReview(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "user_id is required",
			},
		})
		return
	}

	prs, err := c.userService.GetUserReviews(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	pullRequests := make([]gin.H, 0)
	for _, pr := range prs {
		pullRequests = append(pullRequests, gin.H{
			"pull_request_id":   pr.ID,
			"pull_request_name": pr.Name,
			"author_id":         pr.AuthorID,
			"status":            pr.Status,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": pullRequests,
	})
}

// GetStats возвращает статистику PR пользователя
// GET /users/getStats?user_id=<id>
func (c *UserController) GetStats(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "user_id is required",
			},
		})
		return
	}

	stats, err := c.userService.GetUserPRStats(userID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "user not found",
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

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":    stats.UserID,
		"username":   stats.Username,
		"total_prs":  stats.TotalPRs,
		"open_prs":   stats.OpenPRs,
		"merged_prs": stats.MergedPRs,
	})
}
