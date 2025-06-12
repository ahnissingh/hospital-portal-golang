package services_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hospital-project/internal/models"
	"hospital-project/internal/services"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func TestAuthService_Login_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create auth service with mock repository
	authService := services.NewAuthService(mockRepo)

	// Hash the password we'll use in the test
	hashedPassword, err := authService.HashPassword("password123")
	assert.NoError(t, err)

	// Create test user with the hashed password
	user := &models.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
		Role:         models.RoleReceptionist,
	}

	// Set up expectations
	mockRepo.On("FindByUsername", "testuser").Return(user, nil)

	// Create login request
	loginRequest := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	// Call the method being tested
	response, err := authService.Login(loginRequest)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, user.Username, response.User.Username)
	assert.Equal(t, user.Role, response.User.Role)

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create auth service with mock repository
	authService := services.NewAuthService(mockRepo)

	// Hash the password we'll use in the test
	hashedPassword, err := authService.HashPassword("password123")
	assert.NoError(t, err)

	// Create test user with the hashed password
	user := &models.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
		Role:         models.RoleReceptionist,
	}

	// Set up expectations
	mockRepo.On("FindByUsername", "testuser").Return(user, nil)

	// Create login request with wrong password
	loginRequest := models.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	// Call the method being tested
	response, err := authService.Login(loginRequest)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid credentials", err.Error())

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Set up expectations
	mockRepo.On("FindByUsername", "nonexistentuser").Return(nil, errors.New("user not found"))

	// Create auth service with mock repository
	authService := services.NewAuthService(mockRepo)

	// Create login request with non-existent user
	loginRequest := models.LoginRequest{
		Username: "nonexistentuser",
		Password: "password123",
	}

	// Call the method being tested
	response, err := authService.Login(loginRequest)

	// Assert expectations
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "invalid credentials", err.Error())

	// Verify that the mock was called as expected
	mockRepo.AssertExpectations(t)
}

func TestAuthService_HashPassword(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockUserRepository)

	// Create auth service with mock repository
	authService := services.NewAuthService(mockRepo)

	// Call the method being tested
	hashedPassword, err := authService.HashPassword("password123")

	// Assert expectations
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Verify that the hash is valid
	err = authService.VerifyPassword(hashedPassword, "password123")
	assert.NoError(t, err)

	// Verify that the hash is invalid for wrong password
	err = authService.VerifyPassword(hashedPassword, "wrongpassword")
	assert.Error(t, err)
}
