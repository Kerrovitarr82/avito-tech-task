package main

import (
	"avito-tech-task/internal/controllers"
	"avito-tech-task/internal/database"
	"avito-tech-task/internal/services"
	"avito-tech-task/internal/transport"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func initConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	initConfig()

	db := database.ConnectToDB()
	if db == nil {
		log.Fatal("Already exited because of DB, but just in case...")
	}
	database.MigrateDB(db)

	teamService := services.NewTeamService(db)
	userService := services.NewUserService(db)
	prService := services.NewPRService(db)

	database.SeedDB(db, prService)

	teamController := controllers.NewTeamController(teamService)
	userController := controllers.NewUserController(userService)
	prController := controllers.NewPRController(prService)

	router := gin.Default()

	transport.SetupRoutes(router, teamController, userController, prController)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
