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

	t.Run("Valid Token", func(t *testing.T) {
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
	})

	t.Run("Wrong Secret", func(t *testing.T) {
		token, err := GenerateToken(userID, secret, expiration)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		_, err = ValidateToken(token, "wrong-secret")
		if err == nil {
			t.Error("expected error for wrong secret, got nil")
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		// Create a token that is already expired
		token, err := GenerateToken(userID, secret, -time.Hour)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		_, err = ValidateToken(token, secret)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	t.Run("Tampered Token", func(t *testing.T) {
		token, err := GenerateToken(userID, secret, expiration)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		// Tamper with the token
		tamperedToken := token + "tamper"

		_, err = ValidateToken(tamperedToken, secret)
		if err == nil {
			t.Error("expected error for tampered token, got nil")
		}
	})
}
