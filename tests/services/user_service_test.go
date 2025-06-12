package services_test

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hospital-project/internal/models"
	"hospital-project/internal/services"
)

// MockAuthService is a mock implementation of the AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(request models.LoginRequest) (*models.LoginResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockAuthService) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(token *jwt.Token) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_Create_Success(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// Set up expectations
	mockUserRepo.On("FindByUsername", "testuser").Return(nil, errors.New("not found"))
	mockAuthService.On("HashPassword", "password123").Return("hashed_password", nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Create user service with mock repositories
	userService := services.NewUserService(mockUserRepo, mockAuthService)

	// Call the method being tested
	user, err := userService.Create("testuser", "password123", models.RoleReceptionist)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hashed_password", user.PasswordHash)
	assert.Equal(t, models.RoleReceptionist, user.Role)

	// Verify that the mocks were called as expected
	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserService_Create_UsernameExists(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// Create existing user
	existingUser := &models.User{
		Username: "testuser",
		Role:     models.RoleReceptionist,
	}

	// Set up expectations
	mockUserRepo.On("FindByUsername", "testuser").Return(existingUser, nil)

	// Create user service with mock repositories
	userService := services.NewUserService(mockUserRepo, mockAuthService)

	// Call the method being tested
	user, err := userService.Create("testuser", "password123", models.RoleReceptionist)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "username already exists", err.Error())

	// Verify that the mocks were called as expected
	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserService_GetByID_Success(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// Create test user
	user := &models.User{
		Username: "testuser",
		Role:     models.RoleReceptionist,
	}

	// Set up expectations
	mockUserRepo.On("FindByID", uint(1)).Return(user, nil)

	// Create user service with mock repositories
	userService := services.NewUserService(mockUserRepo, mockAuthService)

	// Call the method being tested
	result, err := userService.GetByID(1)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Role, result.Role)

	// Verify that the mocks were called as expected
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// Set up expectations
	mockUserRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	// Create user service with mock repositories
	userService := services.NewUserService(mockUserRepo, mockAuthService)

	// Call the method being tested
	result, err := userService.GetByID(1)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "not found", err.Error())

	// Verify that the mocks were called as expected
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_List_Success(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(MockUserRepository)
	mockAuthService := new(MockAuthService)

	// Create test users
	users := []models.User{
		{
			Username: "user1",
			Role:     models.RoleReceptionist,
		},
		{
			Username: "user2",
			Role:     models.RoleDoctor,
		},
	}

	// Set up expectations
	mockUserRepo.On("List").Return(users, nil)

	// Create user service with mock repositories
	userService := services.NewUserService(mockUserRepo, mockAuthService)

	// Call the method being tested
	result, err := userService.List()

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].Username, result[0].Username)
	assert.Equal(t, users[1].Username, result[1].Username)

	// Verify that the mocks were called as expected
	mockUserRepo.AssertExpectations(t)
}
