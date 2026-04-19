package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"prompt-management/internal/config"
	"prompt-management/internal/domain"
	"prompt-management/pkg/auth"
)

type AuthService struct {
	config *config.Config
	store  domain.UserStore
}

func NewAuthService(cfg *config.Config, store domain.UserStore) *AuthService {
	return &AuthService{
		config: cfg,
		store:  store,
	}
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, email, username, fullName, password string) (*domain.User, error) {
	// 1. Validation basics
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.ToLower(strings.TrimSpace(username))

	// 2. Check if user already exists (Email or Username)
	_, err := s.store.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrConflict // Email already in use
	}
	
	_, err = s.store.GetByUsername(ctx, username)
	if err == nil {
		return nil, domain.ErrConflict // Username already in use
	}

	// 3. Hash Password
	hashed, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 4. Create User
	user := &domain.User{
		Email:        email,
		Username:     username,
		FullName:     fullName,
		PasswordHash: hashed,
	}

	err = s.store.Insert(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login verifies credentials and returns a JWT.
func (s *AuthService) Login(ctx context.Context, identifier, password string) (string, error) {
	var user *domain.User
	var err error

	// 1. Try to find user by email or username
	if strings.Contains(identifier, "@") {
		user, err = s.store.GetByEmail(ctx, strings.ToLower(identifier))
	} else {
		user, err = s.store.GetByUsername(ctx, strings.ToLower(identifier))
	}

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrUnauthorized
		}
		return "", err
	}

	// 2. Verify Password
	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", domain.ErrUnauthorized
	}

	// 3. Generate JWT
	token, err := auth.GenerateToken(user.ID, s.config.JWTSecret, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// RefreshToken generates a new JWT token for an already authenticated user, extending their session.
func (s *AuthService) RefreshToken(ctx context.Context, userID string) (string, error) {
	// A strictly stateless token refresh mechanism expects the identity to be firmly embedded contextually from upstream authorization middleware.
	// As we do not track refresh tokens in PostgreSQL, the action natively extends by re-signing a fresh payload matching the originating UserID.
	token, err := auth.GenerateToken(userID, s.config.JWTSecret, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}
	
	return token, nil
}
