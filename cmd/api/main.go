package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"hospital-project/internal/config"
	"hospital-project/internal/controllers"
	"hospital-project/internal/middleware"
	"hospital-project/internal/repositories"
	"hospital-project/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	patientRepo := repositories.NewPatientRepository(db)

	// Initialize services
	authService := services.NewAuthService()
	userService := services.NewUserService(userRepo, authService)
	patientService := services.NewPatientService(patientRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userService, authService)

	// Initialize router
	router := gin.Default()

	// Initialize controllers
	authController := controllers.NewAuthController(userService, authService)
	adminController := controllers.NewAdminController(userService, patientService, authMiddleware)

	// Register routes
	authController.RegisterRoutes(router)
	adminController.RegisterRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
