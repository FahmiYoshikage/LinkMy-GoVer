package repository

import (
	"context"
	"errors"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository struct {
	db *pgxpool.Pool
}

func NewLinkRepository(db *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{db: db}
}

// Create creates a new link
func (r *LinkRepository) Create(ctx context.Context, link *models.Link) error {
	// Get next position
	var maxPos int
	r.db.QueryRow(ctx, 
		"SELECT COALESCE(MAX(position), 0) FROM links WHERE profile_id = $1",
		link.ProfileID,
	).Scan(&maxPos)
	
	if link.Position == 0 {
		link.Position = maxPos + 1
	}

	query := `
		INSERT INTO links (profile_id, category_id, title, url, icon, position, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		link.ProfileID, link.CategoryID, link.Title, link.URL,
		link.Icon, link.Position, link.IsActive,
	).Scan(&link.ID, &link.CreatedAt)
}

// GetByID retrieves a link by ID
func (r *LinkRepository) GetByID(ctx context.Context, id int) (*models.Link, error) {
	query := `
		SELECT id, profile_id, category_id, title, url, icon, position, clicks, is_active, created_at, updated_at
		FROM links WHERE id = $1
	`
	link := &models.Link{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&link.ID, &link.ProfileID, &link.CategoryID, &link.Title, &link.URL,
		&link.Icon, &link.Position, &link.Clicks, &link.IsActive,
		&link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return link, nil
}

// GetByProfileID retrieves all links for a profile
func (r *LinkRepository) GetByProfileID(ctx context.Context, profileID int, activeOnly bool) ([]models.Link, error) {
	query := `
		SELECT id, profile_id, category_id, title, url, icon, position, clicks, is_active, created_at, updated_at
		FROM links WHERE profile_id = $1
	`
	if activeOnly {
		query += " AND is_active = true"
	}
	query += " ORDER BY position ASC"

	rows, err := r.db.Query(ctx, query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		err := rows.Scan(
			&l.ID, &l.ProfileID, &l.CategoryID, &l.Title, &l.URL,
			&l.Icon, &l.Position, &l.Clicks, &l.IsActive,
			&l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}

// Update updates a link
func (r *LinkRepository) Update(ctx context.Context, link *models.Link) error {
	query := `
		UPDATE links SET category_id = $1, title = $2, url = $3, icon = $4, 
			   position = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`
	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		link.CategoryID, link.Title, link.URL, link.Icon,
		link.Position, link.IsActive, now, link.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete deletes a link
func (r *LinkRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, "DELETE FROM links WHERE id = $1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// IncrementClicks increments the click counter
func (r *LinkRepository) IncrementClicks(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "UPDATE links SET clicks = clicks + 1 WHERE id = $1", id)
	return err
}

// RecordClick records a click event for analytics
func (r *LinkRepository) RecordClick(ctx context.Context, click *models.Click) error {
	query := `
		INSERT INTO clicks (link_id, ip, country, city, user_agent, referrer)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, clicked_at
	`
	return r.db.QueryRow(ctx, query,
		click.LinkID, click.IP, click.Country, click.City, click.UserAgent, click.Referrer,
	).Scan(&click.ID, &click.ClickedAt)
}

// Reorder updates positions for multiple links
func (r *LinkRepository) Reorder(ctx context.Context, positions []models.LinkPosition) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, pos := range positions {
		_, err := tx.Exec(ctx, 
			"UPDATE links SET position = $1, updated_at = $2 WHERE id = $3",
			pos.Position, time.Now(), pos.ID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetProfileOwner gets the user ID that owns this link
func (r *LinkRepository) GetProfileOwner(ctx context.Context, linkID int) (int, error) {
	var userID int
	err := r.db.QueryRow(ctx, `
		SELECT p.user_id FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1
	`, linkID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return userID, nil
}
