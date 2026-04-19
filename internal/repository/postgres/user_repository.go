package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"prompt-management/internal/domain"
)

type userRepo struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(db *pgxpool.Pool) domain.UserStore {
	return &userRepo{db: db}
}

func (r *userRepo) Insert(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, username, full_name, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	args := []interface{}{
		user.Email,
		user.Username,
		user.FullName,
		user.PasswordHash,
	}

	err := r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("could not insert user: %w", err)
	}

	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	return r.findOne(ctx, query, email)
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1 AND deleted_at IS NULL`

	return r.findOne(ctx, query, username)
}

func (r *userRepo) findOne(ctx context.Context, query string, arg interface{}) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, query, arg).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.FullName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("could not find user: %w", err)
	}

	return &user, nil
}
