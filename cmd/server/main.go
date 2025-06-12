package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"hospital-project/internal/config"
	"hospital-project/internal/controllers"
	"hospital-project/internal/middleware"
	"hospital-project/internal/repositories"
	"hospital-project/internal/services"
)

// @title Hospital Management System API
// @version 1.0
// @description API for hospital management system with receptionist and doctor portals
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize database
	dbConfig := config.NewDatabaseConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate database
	err = config.MigrateDB(db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	patientRepo := repositories.NewPatientRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo, authService)
	patientService := services.NewPatientService(patientRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService, authMiddleware)
	patientController := controllers.NewPatientController(patientService, authMiddleware)

	// Initialize router
	router := gin.Default()

	// Register routes
	authController.RegisterRoutes(router)
	userController.RegisterRoutes(router)
	patientController.RegisterRoutes(router)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("Server running on port %s\n", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
