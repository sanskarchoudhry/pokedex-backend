package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sanskarchoudhry/pokedex-backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

// CreateUser inserts a new user into the database
func (r *postgresUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	// QueryRowContext is used because we expect 1 row back (the ID)
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password, // Note: This should be the HASHED password
		time.Now(),
		time.Now(),
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// GetUserByEmail fetches a user by their email address
func (r *postgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository error: %w", err)
	}

	return &user, nil
}
