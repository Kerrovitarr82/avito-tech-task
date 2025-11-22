package transport

import (
	"avito-tech-task/internal/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	teamController *controllers.TeamController,
	userController *controllers.UserController,
	prController *controllers.PRController,
) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	team := router.Group("/team")
	{
		team.POST("/add", teamController.AddTeam)
		team.GET("/get", teamController.GetTeam)
		team.POST("/deactivateMembers", teamController.DeactivateMembers)
	}

	users := router.Group("/users")
	{
		users.POST("/setIsActive", userController.SetIsActive)
		users.GET("/getReview", userController.GetReview)
		users.GET("/getStats", userController.GetStats)
	}

	prs := router.Group("/pullRequest")
	{
		prs.POST("/create", prController.CreatePR)
		prs.POST("/merge", prController.MergePR)
		prs.POST("/reassign", prController.ReassignReviewer)
	}

}
