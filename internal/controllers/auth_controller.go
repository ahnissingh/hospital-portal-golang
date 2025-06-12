package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"hospital-project/internal/models"
	"hospital-project/internal/services"
)

// AuthController handles authentication requests
type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
	}
}

// @Summary Login
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var request models.LoginRequest

	// Bind request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if request.Username == "" || request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Login
	response, err := c.authService.Login(request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if the client is a web browser or Postman
	userAgent := ctx.Request.Header.Get("User-Agent")
	if userAgent != "" && (strings.Contains(userAgent, "PostmanRuntime") ||
		strings.Contains(userAgent, "Mozilla") ||
		strings.Contains(userAgent, "Chrome") ||
		strings.Contains(userAgent, "Safari") ||
		strings.Contains(userAgent, "Firefox") ||
		strings.Contains(userAgent, "Edge")) {
		// Set HTTP-only cookie
		ctx.SetCookie(
			"jwt_token",
			response.Token,
			3600*24, // 24 hours
			"/",
			"",
			false, // secure should be true in production with HTTPS
			true,  // HTTP only
		)
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Register
// @Description Register a new user (doctor or receptionist)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 201 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var request models.RegisterRequest

	// Bind request body
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if request.Username == "" || request.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Validate role
	if request.Role != models.RoleDoctor && request.Role != models.RoleReceptionist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Role must be either doctor or receptionist"})
		return
	}

	// Create user
	user, err := c.userService.Create(request.Username, request.Password, request.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := c.authService.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Create response
	response := &models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	// Check if the client is a web browser or Postman
	userAgent := ctx.Request.Header.Get("User-Agent")
	if userAgent != "" && (strings.Contains(userAgent, "PostmanRuntime") ||
		strings.Contains(userAgent, "Mozilla") ||
		strings.Contains(userAgent, "Chrome") ||
		strings.Contains(userAgent, "Safari") ||
		strings.Contains(userAgent, "Firefox") ||
		strings.Contains(userAgent, "Edge")) {
		// Set HTTP-only cookie
		ctx.SetCookie(
			"jwt_token",
			token,
			3600*24, // 24 hours
			"/",
			"",
			false, // secure should be true in production with HTTPS
			true,  // HTTP only
		)
	}

	ctx.JSON(http.StatusCreated, response)
}

// RegisterRoutes registers the auth routes
func (c *AuthController) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/register", c.Register)
	}
}
