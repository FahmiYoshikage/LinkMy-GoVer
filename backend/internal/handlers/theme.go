package handlers

import (
	"context"
	"errors"

	"github.com/FahmiYoshikage/linkmy-v2/internal/middleware"
	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/FahmiYoshikage/linkmy-v2/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThemeHandler struct {
	themeRepo   *repository.ThemeRepository
	profileRepo *repository.ProfileRepository
}

func NewThemeHandler(db *pgxpool.Pool) *ThemeHandler {
	return &ThemeHandler{
		themeRepo:   repository.NewThemeRepository(db),
		profileRepo: repository.NewProfileRepository(db),
	}
}

// GetTheme returns the theme for a profile
func (h *ThemeHandler) GetTheme(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	profileID, err := c.ParamsInt("profileId")
	if err != nil {
		return ValidationError(c, "Invalid profile ID")
	}

	ctx := context.Background()

	// Check ownership
	belongs, err := h.profileRepo.BelongsToUser(ctx, profileID, userID)
	if err != nil || !belongs {
		return Forbidden(c)
	}

	theme, err := h.themeRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NotFound(c, "Theme")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch theme")
	}

	return SuccessResponse(c, theme)
}

// UpdateTheme updates the theme for a profile
func (h *ThemeHandler) UpdateTheme(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	profileID, err := c.ParamsInt("profileId")
	if err != nil {
		return ValidationError(c, "Invalid profile ID")
	}

	ctx := context.Background()

	// Check ownership
	belongs, err := h.profileRepo.BelongsToUser(ctx, profileID, userID)
	if err != nil || !belongs {
		return Forbidden(c)
	}

	theme, err := h.themeRepo.GetByProfileID(ctx, profileID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NotFound(c, "Theme")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch theme")
	}

	var req models.UpdateThemeRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	// Update fields if provided
	if req.BgType != nil {
		theme.BgType = *req.BgType
	}
	if req.BgValue != nil {
		theme.BgValue = req.BgValue
	}
	if req.ButtonStyle != nil {
		theme.ButtonStyle = *req.ButtonStyle
	}
	if req.ButtonColor != nil {
		theme.ButtonColor = *req.ButtonColor
	}
	if req.TextColor != nil {
		theme.TextColor = *req.TextColor
	}
	if req.Font != nil {
		theme.Font = *req.Font
	}
	if req.Layout != nil {
		theme.Layout = *req.Layout
	}
	if req.ContainerStyle != nil {
		theme.ContainerStyle = *req.ContainerStyle
	}
	if req.EnableAnimations != nil {
		theme.EnableAnimations = *req.EnableAnimations
	}
	if req.EnableGlassEffect != nil {
		theme.EnableGlassEffect = *req.EnableGlassEffect
	}
	if req.ShadowIntensity != nil {
		theme.ShadowIntensity = *req.ShadowIntensity
	}
	if req.BoxedEnabled != nil {
		theme.BoxedEnabled = *req.BoxedEnabled
	}
	if req.BoxedOuterBgType != nil {
		theme.BoxedOuterBgType = req.BoxedOuterBgType
	}
	if req.BoxedOuterBgValue != nil {
		theme.BoxedOuterBgValue = req.BoxedOuterBgValue
	}
	if req.BoxedContainerBg != nil {
		theme.BoxedContainerBg = *req.BoxedContainerBg
	}
	if req.BoxedMaxWidth != nil {
		theme.BoxedMaxWidth = *req.BoxedMaxWidth
	}
	if req.BoxedRadius != nil {
		theme.BoxedRadius = *req.BoxedRadius
	}
	if req.BoxedShadow != nil {
		theme.BoxedShadow = *req.BoxedShadow
	}

	if err := h.themeRepo.Update(ctx, theme); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update theme")
	}

	return SuccessResponse(c, theme)
}
