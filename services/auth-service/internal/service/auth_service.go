package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/XRS0/blog/services/auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	logger    *slog.Logger
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *AuthService) Register(email, username, password string) (*repository.User, string, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, "", fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.userRepo.Create(email, username, string(hashedPassword))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("user registered", "user_id", user.ID, "email", email)
	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*repository.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info("user logged in", "user_id", user.ID, "email", email)
	return user, token, nil
}

func (s *AuthService) ValidateToken(token string) (uint64, error) {
	claims := &jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract user ID from subject
	var userID uint64
	_, err = fmt.Sscanf(claims.Subject, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("invalid token subject")
	}

	return userID, nil
}

func (s *AuthService) GetUserByID(id uint64) (*repository.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) GetUserByEmail(email string) (*repository.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *AuthService) generateToken(userID uint64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
