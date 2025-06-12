package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hospital-project/internal/middleware"
	"hospital-project/internal/models"
	"hospital-project/internal/services"
)

// UserController handles user requests
type UserController struct {
	userService    services.UserService
	authMiddleware *middleware.AuthMiddleware
}

// NewUserController creates a new user controller
func NewUserController(userService services.UserService, authMiddleware *middleware.AuthMiddleware) *UserController {
	return &UserController{
		userService:    userService,
		authMiddleware: authMiddleware,
	}
}

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users [post]
// @Security Bearer
func (c *UserController) CreateUser(ctx *gin.Context) {
	var request struct {
		Username string      `json:"username" binding:"required"`
		Password string      `json:"password" binding:"required,min=6"`
		Role     models.Role `json:"role" binding:"required"`
	}

	// Bind request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create user
	user, err := c.userService.Create(request.Username, request.Password, request.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user.ToResponse())
}

// @Summary Get user by ID
// @Description Get a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [get]
// @Security Bearer
func (c *UserController) GetUser(ctx *gin.Context) {
	// Get ID from path
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user
	user, err := c.userService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user.ToResponse())
}

// @Summary Get current user
// @Description Get the current authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/me [get]
// @Security Bearer
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	// Get user from context
	user, ok := middleware.GetCurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, user.ToResponse())
}

// RegisterRoutes registers the user routes
func (c *UserController) RegisterRoutes(router *gin.Engine) {
	users := router.Group("/api/users")
	users.Use(c.authMiddleware.Authenticate())
	{
		users.POST("", c.CreateUser)
		users.GET("/:id", c.GetUser)
		users.GET("/me", c.GetCurrentUser)
	}
}
