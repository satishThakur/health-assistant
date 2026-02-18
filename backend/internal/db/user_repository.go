package db

import (
	"context"
	"fmt"

	"github.com/satishthakur/health-assistant/backend/internal/models"
)

// UserRepository handles database operations for users.
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

// FindOrCreateUserByGoogleID upserts a user based on their Google ID.
// On conflict it updates email and display_name to reflect the latest Google profile.
func (r *UserRepository) FindOrCreateUserByGoogleID(ctx context.Context, googleID, email, displayName string) (*models.User, error) {
	query := `
		INSERT INTO users (email, google_id, display_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (google_id) DO UPDATE SET
			email        = EXCLUDED.email,
			display_name = EXCLUDED.display_name
		RETURNING id, email, COALESCE(display_name, ''), COALESCE(google_id, ''), created_at
	`

	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, email, googleID, displayName).
		Scan(&user.ID, &user.Email, &user.DisplayName, &user.GoogleID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("find or create user by google ID: %w", err)
	}

	return &user, nil
}

// FindUserByID retrieves a user by their primary key.
func (r *UserRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, email, COALESCE(display_name, ''), COALESCE(google_id, ''), created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.DisplayName, &user.GoogleID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("find user by ID: %w", err)
	}

	return &user, nil
}
