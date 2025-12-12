package repository

import (
	"context"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	// Get next position
	var maxPos int
	r.db.QueryRow(ctx, 
		"SELECT COALESCE(MAX(position), 0) FROM categories WHERE profile_id = $1",
		category.ProfileID,
	).Scan(&maxPos)
	
	if category.Position == 0 {
		category.Position = maxPos + 1
	}

	query := `
		INSERT INTO categories (profile_id, name, icon, color, position, is_expanded)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		category.ProfileID, category.Name, category.Icon, category.Color,
		category.Position, category.IsExpanded,
	).Scan(&category.ID, &category.CreatedAt)
}

// GetByProfileID retrieves all categories for a profile
func (r *CategoryRepository) GetByProfileID(ctx context.Context, profileID int) ([]models.Category, error) {
	query := `
		SELECT id, profile_id, name, icon, color, position, is_expanded, created_at
		FROM categories WHERE profile_id = $1
		ORDER BY position ASC
	`
	rows, err := r.db.Query(ctx, query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(
			&c.ID, &c.ProfileID, &c.Name, &c.Icon, &c.Color,
			&c.Position, &c.IsExpanded, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// Update updates a category
func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories SET name = $1, icon = $2, color = $3, position = $4, is_expanded = $5
		WHERE id = $6
	`
	result, err := r.db.Exec(ctx, query,
		category.Name, category.Icon, category.Color, category.Position, category.IsExpanded, category.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete deletes a category (links will have category_id set to NULL)
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SessionRepository for managing JWT sessions
type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, ip, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		session.UserID, session.RefreshToken, session.IP, session.UserAgent, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
}

// GetByRefreshToken retrieves a session by refresh token
func (r *SessionRepository) GetByRefreshToken(ctx context.Context, token string) (*models.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, ip, user_agent, expires_at, created_at
		FROM sessions WHERE refresh_token = $1 AND expires_at > $2
	`
	session := &models.Session{}
	err := r.db.QueryRow(ctx, query, token, time.Now()).Scan(
		&session.ID, &session.UserID, &session.RefreshToken,
		&session.IP, &session.UserAgent, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Delete deletes a session
func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sessions WHERE id = $1", id)
	return err
}

// DeleteByUserID deletes all sessions for a user
func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sessions WHERE user_id = $1", userID)
	return err
}

// DeleteExpired removes expired sessions
func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.Exec(ctx, "DELETE FROM sessions WHERE expires_at < $1", time.Now())
	return err
}
