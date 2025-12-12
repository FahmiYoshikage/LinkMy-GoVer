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

type LinkHandler struct {
	linkRepo    *repository.LinkRepository
	profileRepo *repository.ProfileRepository
}

func NewLinkHandler(db *pgxpool.Pool) *LinkHandler {
	return &LinkHandler{
		linkRepo:    repository.NewLinkRepository(db),
		profileRepo: repository.NewProfileRepository(db),
	}
}

// GetLinks returns all links for a profile
func (h *LinkHandler) GetLinks(c *fiber.Ctx) error {
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

	links, err := h.linkRepo.GetByProfileID(ctx, profileID, false)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch links")
	}

	return SuccessResponse(c, links)
}

// CreateLink creates a new link
func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
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

	var req models.CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	if req.Title == "" || req.URL == "" {
		return ValidationError(c, "Title and URL are required")
	}

	icon := "bi-link-45deg"
	if req.Icon != "" {
		icon = req.Icon
	}

	link := &models.Link{
		ProfileID:  profileID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		URL:        req.URL,
		Icon:       icon,
		IsActive:   true,
	}

	if req.Position != nil {
		link.Position = *req.Position
	}

	if err := h.linkRepo.Create(ctx, link); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create link")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    link,
	})
}

// UpdateLink updates a link
func (h *LinkHandler) UpdateLink(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	linkID, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid link ID")
	}

	ctx := context.Background()

	// Check ownership
	ownerID, err := h.linkRepo.GetProfileOwner(ctx, linkID)
	if err != nil {
		return NotFound(c, "Link")
	}
	if ownerID != userID {
		return Forbidden(c)
	}

	link, err := h.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		return NotFound(c, "Link")
	}

	var req models.UpdateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	// Update fields if provided
	if req.Title != nil {
		link.Title = *req.Title
	}
	if req.URL != nil {
		link.URL = *req.URL
	}
	if req.Icon != nil {
		link.Icon = *req.Icon
	}
	if req.CategoryID != nil {
		link.CategoryID = req.CategoryID
	}
	if req.Position != nil {
		link.Position = *req.Position
	}
	if req.IsActive != nil {
		link.IsActive = *req.IsActive
	}

	if err := h.linkRepo.Update(ctx, link); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update link")
	}

	return SuccessResponse(c, link)
}

// DeleteLink deletes a link
func (h *LinkHandler) DeleteLink(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	linkID, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid link ID")
	}

	ctx := context.Background()

	// Check ownership
	ownerID, err := h.linkRepo.GetProfileOwner(ctx, linkID)
	if err != nil {
		return NotFound(c, "Link")
	}
	if ownerID != userID {
		return Forbidden(c)
	}

	if err := h.linkRepo.Delete(ctx, linkID); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete link")
	}

	return SuccessResponse(c, fiber.Map{"message": "Link deleted"})
}

// ReorderLinks updates link positions
func (h *LinkHandler) ReorderLinks(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req models.ReorderLinksRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	ctx := context.Background()

	// Verify ownership of all links
	for _, pos := range req.Links {
		ownerID, err := h.linkRepo.GetProfileOwner(ctx, pos.ID)
		if err != nil || ownerID != userID {
			return Forbidden(c)
		}
	}

	if err := h.linkRepo.Reorder(ctx, req.Links); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to reorder links")
	}

	return SuccessResponse(c, fiber.Map{"message": "Links reordered"})
}

// TrackClick records a link click (public endpoint)
func (h *LinkHandler) TrackClick(c *fiber.Ctx) error {
	linkID, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid link ID")
	}

	var req models.TrackClickRequest
	c.BodyParser(&req) // Optional body

	ctx := context.Background()

	// Verify link exists
	link, err := h.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return NotFound(c, "Link")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Increment click counter
	h.linkRepo.IncrementClicks(ctx, linkID)

	// Record click for analytics
	ip := c.IP()
	userAgent := c.Get("User-Agent")
	click := &models.Click{
		LinkID:    linkID,
		IP:        &ip,
		UserAgent: &userAgent,
		Referrer:  req.Referrer,
		// Country and city would be looked up via IP geolocation service
	}
	h.linkRepo.RecordClick(ctx, click)

	return SuccessResponse(c, fiber.Map{
		"url": link.URL,
	})
}
