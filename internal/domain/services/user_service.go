package services

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lealencar/financial-api/internal/domain/models"
	"github.com/lealencar/financial-api/internal/domain/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic related to users
type UserService struct {
	userRepo      *repositories.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

// NewUserService creates a new user service
func NewUserService(userRepo *repositories.UserRepository, jwtSecret string, jwtExpiration time.Duration) *UserService {
	return &UserService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

// RegisterUserInput contains data needed to register a new user
type RegisterUserInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// RegisterUserOutput contains the result of a successful registration
type RegisterUserOutput struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

// RegisterUser creates a new user account
func (s *UserService) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	// Validate email format
	if !isValidEmail(input.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password strength
	if !isStrongPassword(input.Password) {
		return nil, errors.New("password must be at least 8 characters and include uppercase, lowercase, number, and special character")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Create new user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Email:     input.Email,
		Password:  string(hashedPassword),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Return user data and token
	return &RegisterUserOutput{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		Token:     token,
	}, nil
}

// generateToken creates a new JWT token for a user
func (s *UserService) generateToken(userID uint) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(s.jwtExpiration)

	// Create claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
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
