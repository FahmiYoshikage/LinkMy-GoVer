package repository

import (
	"context"
	"errors"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicate = errors.New("already exists")

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, is_verified, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.IsVerified, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt)
	
	if err != nil {
		if isDuplicateError(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_verified, is_active, is_admin, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = true
	`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsVerified, &user.IsActive, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_verified, is_active, is_admin, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsVerified, &user.IsActive, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_verified, is_active, is_admin, created_at, updated_at
		FROM users WHERE username = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsVerified, &user.IsActive, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET username = $1, email = $2, is_verified = $3, updated_at = $4
		WHERE id = $5
	`
	now := time.Now()
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.IsVerified, now, user.ID)
	return err
}

// UpdatePassword updates user password
func (r *UserRepository) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, passwordHash, time.Now(), id)
	return err
}

// ExistsEmail checks if email exists
func (r *UserRepository) ExistsEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

// ExistsUsername checks if username exists
func (r *UserRepository) ExistsUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	return exists, err
}

// Helper to check for duplicate key errors
func isDuplicateError(err error) bool {
	return err != nil && (
		contains(err.Error(), "duplicate key") ||
		contains(err.Error(), "unique constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
