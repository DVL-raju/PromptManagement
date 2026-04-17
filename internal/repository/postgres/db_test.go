package postgres

import (
	"testing"
	"prompt-management/internal/config"
)

func TestNewPoolConfig(t *testing.T) {
	// We can't easily test pgxpool.NewWithConfig without a real DB or heavy mocking,
	// but we can verify our configuration logic by ensuring it parses correctly.
	
	t.Run("Invalid Database URL", func(t *testing.T) {
		badCfg := &config.Config{DatabaseURL: "not-a-url"}
		_, err := NewPool(badCfg)
		if err == nil {
			t.Error("expected error for invalid database URL, got nil")
		}
	})

	t.Run("Valid URL Syntax Parsing", func(t *testing.T) {
		validCfg := &config.Config{
			DatabaseURL: "postgres://user:pass@localhost:5432/db",
			DBMaxConns:  "10",
		}
		// NewPool doesn't connect immediately in this setup (background init)
		// but it parses the config.
		p, err := NewPool(validCfg)
		if err != nil {
			t.Errorf("unexpected error for valid URL syntax: %v", err)
		}
		if p != nil {
			p.Close()
		}
	})
}

// Note: Further integration testing would involve a real postgres instance.
