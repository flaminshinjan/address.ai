package service

import (
	"errors"
	"time"

	"github.com/yourusername/hotel-management/pkg/common/auth"
	"github.com/yourusername/hotel-management/services/user/internal/model"
	"github.com/yourusername/hotel-management/services/user/internal/repository"
)

// UserService handles business logic for users
type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

// NewUserService creates a new UserService
func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Register registers a new user
func (s *UserService) Register(reg model.UserRegistration) (model.UserResponse, string, error) {
	// Check if username already exists
	_, err := s.repo.GetByUsername(reg.Username)
	if err == nil {
		return model.UserResponse{}, "", errors.New("username already exists")
	}

	// Check if email already exists
	_, err = s.repo.GetByEmail(reg.Email)
	if err == nil {
		return model.UserResponse{}, "", errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := model.HashPassword(reg.Password)
	if err != nil {
		return model.UserResponse{}, "", err
	}

	// Create user
	user := model.User{
		Username:  reg.Username,
		Email:     reg.Email,
		Password:  hashedPassword,
		FirstName: reg.FirstName,
		LastName:  reg.LastName,
		Role:      "guest", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user, err = s.repo.Create(user)
	if err != nil {
		return model.UserResponse{}, "", err
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return model.UserResponse{}, "", err
	}

	return user.ToResponse(), token, nil
}

// Login logs in a user
func (s *UserService) Login(login model.UserLogin) (model.UserResponse, string, error) {
	// Get user by username
	user, err := s.repo.GetByUsername(login.Username)
	if err != nil {
		return model.UserResponse{}, "", errors.New("invalid username or password")
	}

	// Check password
	if !model.CheckPassword(login.Password, user.Password) {
		return model.UserResponse{}, "", errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return model.UserResponse{}, "", err
	}

	return user.ToResponse(), token, nil
}

// GetByID gets a user by ID
func (s *UserService) GetByID(id string) (model.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return model.UserResponse{}, err
	}

	return user.ToResponse(), nil
}

// Update updates a user
func (s *UserService) Update(id string, user model.User) (model.UserResponse, error) {
	// Get existing user
	existingUser, err := s.repo.GetByID(id)
	if err != nil {
		return model.UserResponse{}, err
	}

	// Update fields
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Role = user.Role
	existingUser.UpdatedAt = time.Now()

	// Save user
	updatedUser, err := s.repo.Update(existingUser)
	if err != nil {
		return model.UserResponse{}, err
	}

	return updatedUser.ToResponse(), nil
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(id, currentPassword, newPassword string) error {
	// Get user
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Check current password
	if !model.CheckPassword(currentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := model.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	return s.repo.UpdatePassword(id, hashedPassword)
}

// Delete deletes a user
func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}

// List lists all users
func (s *UserService) List(limit, offset int) ([]model.UserResponse, error) {
	users, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return userResponses, nil
}
