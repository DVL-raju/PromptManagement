package auth

import (
	"testing"
	"time"
)

func TestPasswordHashing(t *testing.T) {
	password := "super-secret-password"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == password {
		t.Errorf("expected hash to be different from original password")
	}

	if !CheckPasswordHash(password, hash) {
		t.Errorf("expected password to match hash")
	}

	if CheckPasswordHash("wrong-password", hash) {
		t.Errorf("expected wrong password to not match hash")
	}
}

func TestJWTToken(t *testing.T) {
	secret := "test-secret"
	userID := "user-123"
	expiration := time.Hour

	token, err := GenerateToken(userID, secret, expiration)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %s, got %s", userID, claims.UserID)
	}

	_, err = ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}
