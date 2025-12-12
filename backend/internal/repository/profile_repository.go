package repository

import (
	"context"
	"errors"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Create creates a new profile and its default theme
func (r *ProfileRepository) Create(ctx context.Context, profile *models.Profile) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert profile
	query := `
		INSERT INTO profiles (user_id, slug, name, title, bio, avatar, is_active, display_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	err = tx.QueryRow(ctx, query,
		profile.UserID, profile.Slug, profile.Name, profile.Title, profile.Bio,
		profile.Avatar, profile.IsActive, profile.DisplayOrder,
	).Scan(&profile.ID, &profile.CreatedAt)
	
	if err != nil {
		if isDuplicateError(err) {
			return ErrDuplicate
		}
		return err
	}

	// Create default theme
	themeQuery := `
		INSERT INTO themes (profile_id, bg_type, bg_value, button_style, button_color, text_color, font)
		VALUES ($1, 'gradient', 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', 'rounded', '#667eea', '#333333', 'Inter')
	`
	_, err = tx.Exec(ctx, themeQuery, profile.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetByID retrieves a profile by ID
func (r *ProfileRepository) GetByID(ctx context.Context, id int) (*models.Profile, error) {
	query := `
		SELECT id, user_id, slug, name, title, bio, avatar, is_active, display_order, created_at, updated_at
		FROM profiles WHERE id = $1
	`
	profile := &models.Profile{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&profile.ID, &profile.UserID, &profile.Slug, &profile.Name, &profile.Title,
		&profile.Bio, &profile.Avatar, &profile.IsActive, &profile.DisplayOrder,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return profile, nil
}

// GetBySlug retrieves a profile by slug (for public view)
func (r *ProfileRepository) GetBySlug(ctx context.Context, slug string) (*models.Profile, error) {
	query := `
		SELECT id, user_id, slug, name, title, bio, avatar, is_active, display_order, created_at, updated_at
		FROM profiles WHERE slug = $1 AND is_active = true
	`
	profile := &models.Profile{}
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&profile.ID, &profile.UserID, &profile.Slug, &profile.Name, &profile.Title,
		&profile.Bio, &profile.Avatar, &profile.IsActive, &profile.DisplayOrder,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return profile, nil
}

// GetByUserID retrieves all profiles for a user
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID int) ([]models.ProfileWithStats, error) {
	query := `
		SELECT p.id, p.user_id, p.slug, p.name, p.title, p.bio, p.avatar, 
			   p.is_active, p.display_order, p.created_at, p.updated_at,
			   COUNT(DISTINCT l.id) as link_count,
			   COALESCE(SUM(l.clicks), 0) as total_clicks
		FROM profiles p
		LEFT JOIN links l ON l.profile_id = p.id AND l.is_active = true
		WHERE p.user_id = $1
		GROUP BY p.id
		ORDER BY p.display_order ASC, p.created_at ASC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.ProfileWithStats
	for rows.Next() {
		var p models.ProfileWithStats
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Slug, &p.Name, &p.Title, &p.Bio, &p.Avatar,
			&p.IsActive, &p.DisplayOrder, &p.CreatedAt, &p.UpdatedAt,
			&p.LinkCount, &p.TotalClicks,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// Update updates a profile
func (r *ProfileRepository) Update(ctx context.Context, profile *models.Profile) error {
	query := `
		UPDATE profiles SET slug = $1, name = $2, title = $3, bio = $4, 
			   avatar = $5, is_active = $6, display_order = $7, updated_at = $8
		WHERE id = $9
	`
	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		profile.Slug, profile.Name, profile.Title, profile.Bio,
		profile.Avatar, profile.IsActive, profile.DisplayOrder, now, profile.ID,
	)
	if err != nil {
		if isDuplicateError(err) {
			return ErrDuplicate
		}
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete deletes a profile (cascade deletes theme, links, categories)
func (r *ProfileRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, "DELETE FROM profiles WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ExistsSlug checks if slug is taken
func (r *ProfileRepository) ExistsSlug(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM profiles WHERE slug = $1)", slug).Scan(&exists)
	return exists, err
}

// BelongsToUser checks if profile belongs to user
func (r *ProfileRepository) BelongsToUser(ctx context.Context, profileID, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, 
		"SELECT EXISTS(SELECT 1 FROM profiles WHERE id = $1 AND user_id = $2)", 
		profileID, userID,
	).Scan(&exists)
	return exists, err
}
