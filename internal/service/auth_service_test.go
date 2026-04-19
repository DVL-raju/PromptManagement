package service

import (
	"context"
	"testing"

	"prompt-management/internal/config"
	"prompt-management/internal/domain"
)

// MockStore is a simple in-memory implementation of UserStore for testing.
type MockStore struct {
	users map[string]*domain.User
}

func (m *MockStore) Insert(ctx context.Context, user *domain.User) error {
	user.ID = "generated-uuid"
	m.users[user.Email] = user
	return nil
}

func (m *MockStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (m *MockStore) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func TestAuthServiceRegisterAndLogin(t *testing.T) {
	cfg := &config.Config{JWTSecret: "secret"}
	store := &MockStore{users: make(map[string]*domain.User)}
	svc := NewAuthService(cfg, store)
	ctx := context.Background()

	email := "test@example.com"
	_, username, fullName, password := email, "testuser", "Test User", "password123"

	t.Run("Successful Registration", func(t *testing.T) {
		_, err := svc.Register(ctx, email, username, fullName, password)
		if err != nil {
			t.Fatalf("failed to register: %v", err)
		}
	})

	t.Run("Duplicate Registration Email", func(t *testing.T) {
		_, err := svc.Register(ctx, email, "otheruser", fullName, password)
		if err == nil {
			t.Fatal("expected conflict error for existing email, got nil")
		}
	})

	t.Run("Successful Login - Email", func(t *testing.T) {
		token, err := svc.Login(ctx, email, password)
		if err != nil {
			t.Fatalf("failed to login: %v", err)
		}
		if token == "" {
			t.Fatal("expected token, got empty string")
		}
	})

	t.Run("Successful Login - Username", func(t *testing.T) {
		token, err := svc.Login(ctx, username, password)
		if err != nil {
			t.Fatalf("failed to login with username: %v", err)
		}
		if token == "" {
			t.Error("expected token, got empty string")
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		_, err := svc.Login(ctx, email, "wrongpassword")
		if err == nil {
			t.Fatal("expected unauthorized error, got nil")
		}
	})
}
