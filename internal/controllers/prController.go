package controllers

import (
	"avito-tech-task/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PRController struct {
	prService *services.PRService
}

func NewPRController(prService *services.PRService) *PRController {
	return &PRController{
		prService: prService,
	}
}

// CreatePR создает PR и автоматически назначает до 2 ревьюверов
// POST /pullRequest/create
func (c *PRController) CreatePR(ctx *gin.Context) {
	var req struct {
		PullRequestID   string `json:"pull_request_id" binding:"required"`
		PullRequestName string `json:"pull_request_name" binding:"required"`
		AuthorID        string `json:"author_id" binding:"required"`
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

	pr, err := c.prService.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch err.Error() {
		case "PR_EXISTS":
			ctx.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
		case "NOT_FOUND":
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "author or team not found",
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

	reviewers := make([]string, 0)
	for _, reviewer := range pr.Reviewers {
		reviewers = append(reviewers, reviewer.ID)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"pr": gin.H{
			"pull_request_id":    pr.ID,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": reviewers,
			"createdAt":          pr.CreatedAt,
		},
	})
}

// MergePR помечает PR как MERGED (идемпотентная операция)
// POST /pullRequest/merge
func (c *PRController) MergePR(ctx *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id" binding:"required"`
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

	pr, err := c.prService.MergePR(req.PullRequestID)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "PR not found",
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

	reviewers := make([]string, 0)
	for _, reviewer := range pr.Reviewers {
		reviewers = append(reviewers, reviewer.ID)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"pr": gin.H{
			"pull_request_id":    pr.ID,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": reviewers,
			"mergedAt":           pr.MergedAt,
		},
	})
}

// ReassignReviewer переназначает конкретного ревьювера на другого из команды
// POST /pullRequest/reassign
func (c *PRController) ReassignReviewer(ctx *gin.Context) {
	var req struct {
		PullRequestID string `json:"pull_request_id" binding:"required"`
		OldReviewerID string `json:"old_reviewer_id" binding:"required"`
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

	pr, replacedBy, err := c.prService.ReassignReviewer(req.PullRequestID, req.OldReviewerID)
	if err != nil {
		switch err.Error() {
		case "NOT_FOUND":
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "PR or user not found",
				},
			})
		case "PR_MERGED":
			ctx.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "PR_MERGED",
					"message": "cannot reassign on merged PR",
				},
			})
		case "NOT_ASSIGNED":
			ctx.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "NOT_ASSIGNED",
					"message": "reviewer is not assigned to this PR",
				},
			})
		case "NO_CANDIDATE":
			ctx.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "NO_CANDIDATE",
					"message": "no active replacement candidate in team",
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

	reviewers := make([]string, 0)
	for _, reviewer := range pr.Reviewers {
		reviewers = append(reviewers, reviewer.ID)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"pr": gin.H{
			"pull_request_id":    pr.ID,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": reviewers,
		},
		"replaced_by": replacedBy,
	})
}
