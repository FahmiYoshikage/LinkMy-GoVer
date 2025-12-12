package repository

import (
	"context"
	"errors"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThemeRepository struct {
	db *pgxpool.Pool
}

func NewThemeRepository(db *pgxpool.Pool) *ThemeRepository {
	return &ThemeRepository{db: db}
}

// GetByProfileID retrieves theme for a profile
func (r *ThemeRepository) GetByProfileID(ctx context.Context, profileID int) (*models.Theme, error) {
	query := `
		SELECT id, profile_id, bg_type, bg_value, button_style, button_color, text_color, font,
			   layout, container_style, enable_animations, enable_glass_effect, shadow_intensity,
			   boxed_enabled, boxed_outer_bg_type, boxed_outer_bg_value, boxed_container_bg,
			   boxed_max_width, boxed_radius, boxed_shadow, created_at, updated_at
		FROM themes WHERE profile_id = $1
	`
	theme := &models.Theme{}
	err := r.db.QueryRow(ctx, query, profileID).Scan(
		&theme.ID, &theme.ProfileID, &theme.BgType, &theme.BgValue,
		&theme.ButtonStyle, &theme.ButtonColor, &theme.TextColor, &theme.Font,
		&theme.Layout, &theme.ContainerStyle, &theme.EnableAnimations, &theme.EnableGlassEffect,
		&theme.ShadowIntensity, &theme.BoxedEnabled, &theme.BoxedOuterBgType, &theme.BoxedOuterBgValue,
		&theme.BoxedContainerBg, &theme.BoxedMaxWidth, &theme.BoxedRadius, &theme.BoxedShadow,
		&theme.CreatedAt, &theme.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return theme, nil
}

// Update updates a theme
func (r *ThemeRepository) Update(ctx context.Context, theme *models.Theme) error {
	query := `
		UPDATE themes SET 
			bg_type = $1, bg_value = $2, button_style = $3, button_color = $4, text_color = $5,
			font = $6, layout = $7, container_style = $8, enable_animations = $9, enable_glass_effect = $10,
			shadow_intensity = $11, boxed_enabled = $12, boxed_outer_bg_type = $13, boxed_outer_bg_value = $14,
			boxed_container_bg = $15, boxed_max_width = $16, boxed_radius = $17, boxed_shadow = $18,
			updated_at = $19
		WHERE profile_id = $20
	`
	now := time.Now()
	result, err := r.db.Exec(ctx, query,
		theme.BgType, theme.BgValue, theme.ButtonStyle, theme.ButtonColor, theme.TextColor,
		theme.Font, theme.Layout, theme.ContainerStyle, theme.EnableAnimations, theme.EnableGlassEffect,
		theme.ShadowIntensity, theme.BoxedEnabled, theme.BoxedOuterBgType, theme.BoxedOuterBgValue,
		theme.BoxedContainerBg, theme.BoxedMaxWidth, theme.BoxedRadius, theme.BoxedShadow,
		now, theme.ProfileID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Create creates a new theme (usually called when creating profile)
func (r *ThemeRepository) Create(ctx context.Context, theme *models.Theme) error {
	query := `
		INSERT INTO themes (profile_id, bg_type, bg_value, button_style, button_color, text_color, font)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		theme.ProfileID, theme.BgType, theme.BgValue, theme.ButtonStyle,
		theme.ButtonColor, theme.TextColor, theme.Font,
	).Scan(&theme.ID, &theme.CreatedAt)
}
