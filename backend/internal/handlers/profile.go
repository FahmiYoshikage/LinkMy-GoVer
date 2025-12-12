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

type ProfileHandler struct {
	profileRepo  *repository.ProfileRepository
	linkRepo     *repository.LinkRepository
	categoryRepo *repository.CategoryRepository
	themeRepo    *repository.ThemeRepository
	userRepo     *repository.UserRepository
}

func NewProfileHandler(db *pgxpool.Pool) *ProfileHandler {
	return &ProfileHandler{
		profileRepo:  repository.NewProfileRepository(db),
		linkRepo:     repository.NewLinkRepository(db),
		categoryRepo: repository.NewCategoryRepository(db),
		themeRepo:    repository.NewThemeRepository(db),
		userRepo:     repository.NewUserRepository(db),
	}
}

// GetPublicProfile returns a profile for public viewing
func (h *ProfileHandler) GetPublicProfile(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return NotFound(c, "Profile")
	}

	ctx := context.Background()

	// Get profile
	profile, err := h.profileRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NotFound(c, "Profile")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Get theme
	theme, err := h.themeRepo.GetByProfileID(ctx, profile.ID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Get categories
	categories, err := h.categoryRepo.GetByProfileID(ctx, profile.ID)
	if err != nil {
		categories = []models.Category{}
	}

	// Get active links
	links, err := h.linkRepo.GetByProfileID(ctx, profile.ID, true)
	if err != nil {
		links = []models.Link{}
	}

	// Get user verification status
	user, _ := h.userRepo.GetByID(ctx, profile.UserID)
	isVerified := false
	if user != nil {
		isVerified = user.IsVerified
	}

	return SuccessResponse(c, models.PublicProfile{
		Profile:    *profile,
		Theme:      *theme,
		Categories: categories,
		Links:      links,
		IsVerified: isVerified,
	})
}

// GetUserProfiles returns all profiles for the authenticated user
func (h *ProfileHandler) GetUserProfiles(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return Unauthorized(c)
	}

	ctx := context.Background()
	profiles, err := h.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch profiles")
	}

	return SuccessResponse(c, profiles)
}

// GetProfile returns a single profile (owner only)
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	profileID, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid profile ID")
	}

	ctx := context.Background()

	// Check ownership
	belongs, err := h.profileRepo.BelongsToUser(ctx, profileID, userID)
	if err != nil || !belongs {
		return Forbidden(c)
	}

	profile, err := h.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return NotFound(c, "Profile")
	}

	return SuccessResponse(c, profile)
}

// CreateProfile creates a new profile
func (h *ProfileHandler) CreateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return Unauthorized(c)
	}

	var req models.CreateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	if len(req.Slug) < 3 {
		return ValidationError(c, "Slug must be at least 3 characters")
	}

	ctx := context.Background()

	// Check if slug exists
	exists, err := h.profileRepo.ExistsSlug(ctx, req.Slug)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if exists {
		return ErrorResponse(c, fiber.StatusConflict, "Slug already taken")
	}

	profile := &models.Profile{
		UserID:   userID,
		Slug:     req.Slug,
		Name:     req.Name,
		Title:    req.Title,
		Bio:      req.Bio,
		Avatar:   "default-avatar.png",
		IsActive: true,
	}

	if err := h.profileRepo.Create(ctx, profile); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrorResponse(c, fiber.StatusConflict, "Profile already exists")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create profile")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    profile,
	})
}

// UpdateProfile updates a profile
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	profileID, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid profile ID")
	}

	ctx := context.Background()

	// Check ownership
	belongs, err := h.profileRepo.BelongsToUser(ctx, profileID, userID)
	if err != nil || !belongs {
		return Forbidden(c)
	}

	profile, err := h.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return NotFound(c, "Profile")
	}

	var req models.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	// Update fields if provided
	if req.Slug != nil {
		profile.Slug = *req.Slug
	}
	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.Title != nil {
		profile.Title = req.Title
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.Avatar != nil {
		profile.Avatar = *req.Avatar
	}
	if req.IsActive != nil {
		profile.IsActive = *req.IsActive
	}

	if err := h.profileRepo.Update(ctx, profile); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrorResponse(c, fiber.StatusConflict, "Slug already taken")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile")
	}

	return SuccessResponse(c, profile)
}

// DeleteProfile deletes a profile
func (h *ProfileHandler) DeleteProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	profileID, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid profile ID")
	}

	ctx := context.Background()

	// Check ownership
	belongs, err := h.profileRepo.BelongsToUser(ctx, profileID, userID)
	if err != nil || !belongs {
		return Forbidden(c)
	}

	if err := h.profileRepo.Delete(ctx, profileID); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete profile")
	}

	return SuccessResponse(c, fiber.Map{"message": "Profile deleted"})
}
