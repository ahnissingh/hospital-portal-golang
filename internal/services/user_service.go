package services

import (
	"errors"

	"hospital-project/internal/models"
	"hospital-project/internal/repositories"
)

// UserService interface defines methods for user service
type UserService interface {
	Create(username, password string, role models.Role) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List() ([]models.User, error)
}

// userService implements UserService interface
type userService struct {
	userRepo   repositories.UserRepository
	authService AuthService
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository, authService AuthService) UserService {
	return &userService{
		userRepo:   userRepo,
		authService: authService,
	}
}

// Create creates a new user
func (s *userService) Create(username, password string, role models.Role) (*models.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.FindByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := s.authService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
	}

	// Save user to database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID gets a user by ID
func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetByUsername gets a user by username
func (s *userService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.FindByUsername(username)
}

// Update updates a user
func (s *userService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

// Delete deletes a user
func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

// List returns all users
func (s *userService) List() ([]models.User, error) {
	return s.userRepo.List()
}
