package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/shared/messaging/events"
	"github.com/leandroalencar/banco-dados/shared/models"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic related to users
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUserInput contains data needed to create a new user
type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput contains data needed to log in
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse contains the result of a successful login
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// ProcessUserEvent handles different types of user events
func (s *UserService) ProcessUserEvent(ctx context.Context, event *events.UserEvent) error {
	switch event.Action {
	case events.UserActionCreate:
		return s.handleCreate(ctx, event.Data)
	case events.UserActionUpdate:
		return s.handleUpdate(ctx, event.Data)
	case events.UserActionDelete:
		return s.handleDelete(ctx, event.Data)
	default:
		return fmt.Errorf("unknown action type: %s", event.Action)
	}
}

func (s *UserService) handleCreate(ctx context.Context, data events.UserEventData) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: string(hashedPassword),
		// Balance will be set by repository to R$ 1000
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (s *UserService) handleUpdate(ctx context.Context, data events.UserEventData) error {
	// First get the existing user
	user, err := s.userRepo.GetByID(ctx, uint(data.ID))
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	// Update fields if provided
	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Email != "" {
		user.Email = data.Email
	}
	if data.Password != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}
		user.Password = string(hashedPassword)
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (s *UserService) handleDelete(ctx context.Context, data events.UserEventData) error {
	if err := s.userRepo.Delete(ctx, uint(data.ID)); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

// CreateUser creates a new user account
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user - balance will be set by repository to R$ 1000
	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		// Balance will be set by repository to R$ 1000
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login logs in a user
func (s *UserService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// generateJWT creates a new JWT token for a user
func generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-default-secret-key" // You should set this in your .env file
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Helper functions for validation
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func isStrongPassword(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}

	// Check for uppercase, lowercase, number, and special character
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
