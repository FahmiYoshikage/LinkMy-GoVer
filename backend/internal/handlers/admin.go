package handlers

import (
	"context"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminHandler struct {
	db *pgxpool.Pool
}

func NewAdminHandler(db *pgxpool.Pool) *AdminHandler {
	return &AdminHandler{db: db}
}

// Dashboard stats
type AdminStats struct {
	TotalUsers      int `json:"total_users"`
	VerifiedUsers   int `json:"verified_users"`
	TotalProfiles   int `json:"total_profiles"`
	TotalLinks      int `json:"total_links"`
	TotalClicks     int `json:"total_clicks"`
	NewUsersWeek    int `json:"new_users_week"`
	ActiveUsersWeek int `json:"active_users_week"`
}

// GetStats returns dashboard statistics
func (h *AdminHandler) GetStats(c *fiber.Ctx) error {
	ctx := context.Background()
	stats := AdminStats{}

	// Total users
	h.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	
	// Verified users
	h.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE is_verified = true").Scan(&stats.VerifiedUsers)
	
	// Total profiles
	h.db.QueryRow(ctx, "SELECT COUNT(*) FROM profiles").Scan(&stats.TotalProfiles)
	
	// Total links
	h.db.QueryRow(ctx, "SELECT COUNT(*) FROM links").Scan(&stats.TotalLinks)
	
	// Total clicks
	h.db.QueryRow(ctx, "SELECT COALESCE(SUM(clicks), 0) FROM links").Scan(&stats.TotalClicks)
	
	// New users this week
	h.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE created_at > NOW() - INTERVAL '7 days'").Scan(&stats.NewUsersWeek)
	
	// Active users (logged in this week - approximate via sessions)
	h.db.QueryRow(ctx, "SELECT COUNT(DISTINCT user_id) FROM sessions WHERE created_at > NOW() - INTERVAL '7 days'").Scan(&stats.ActiveUsersWeek)

	return SuccessResponse(c, stats)
}

// User list item for admin
type AdminUser struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	IsVerified bool       `json:"is_verified"`
	IsActive   bool       `json:"is_active"`
	IsAdmin    bool       `json:"is_admin"`
	CreatedAt  time.Time  `json:"created_at"`
	ProfileCount int      `json:"profile_count"`
	TotalClicks  int      `json:"total_clicks"`
}

// ListUsers returns all users for admin
func (h *AdminHandler) ListUsers(c *fiber.Ctx) error {
	ctx := context.Background()
	
	search := c.Query("search", "")
	
	query := `
		SELECT u.id, u.username, u.email, u.is_verified, u.is_active, u.is_admin, u.created_at,
			COALESCE((SELECT COUNT(*) FROM profiles WHERE user_id = u.id), 0) as profile_count,
			COALESCE((SELECT SUM(l.clicks) FROM profiles p JOIN links l ON l.profile_id = p.id WHERE p.user_id = u.id), 0) as total_clicks
		FROM users u
		WHERE ($1 = '' OR u.username ILIKE '%' || $1 || '%' OR u.email ILIKE '%' || $1 || '%')
		ORDER BY u.created_at DESC
		LIMIT 100
	`
	
	rows, err := h.db.Query(ctx, query, search)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	defer rows.Close()
	
	users := []AdminUser{}
	for rows.Next() {
		var u AdminUser
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.IsVerified, &u.IsActive, &u.IsAdmin, &u.CreatedAt, &u.ProfileCount, &u.TotalClicks)
		if err != nil {
			continue
		}
		users = append(users, u)
	}
	
	return SuccessResponse(c, users)
}

// UpdateUser updates user status (verify, ban, admin)
func (h *AdminHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid user ID")
	}
	
	var req struct {
		IsVerified *bool `json:"is_verified"`
		IsActive   *bool `json:"is_active"`
		IsAdmin    *bool `json:"is_admin"`
	}
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}
	
	ctx := context.Background()
	
	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argNum := 1
	
	if req.IsVerified != nil {
		updates = append(updates, "is_verified = $"+string(rune('0'+argNum)))
		args = append(args, *req.IsVerified)
		argNum++
	}
	if req.IsActive != nil {
		updates = append(updates, "is_active = $"+string(rune('0'+argNum)))
		args = append(args, *req.IsActive)
		argNum++
	}
	if req.IsAdmin != nil {
		updates = append(updates, "is_admin = $"+string(rune('0'+argNum)))
		args = append(args, *req.IsAdmin)
		argNum++
	}
	
	if len(updates) == 0 {
		return ValidationError(c, "No updates provided")
	}
	
	// Simple approach - update each field individually
	if req.IsVerified != nil {
		h.db.Exec(ctx, "UPDATE users SET is_verified = $1 WHERE id = $2", *req.IsVerified, id)
	}
	if req.IsActive != nil {
		h.db.Exec(ctx, "UPDATE users SET is_active = $1 WHERE id = $2", *req.IsActive, id)
	}
	if req.IsAdmin != nil {
		h.db.Exec(ctx, "UPDATE users SET is_admin = $1 WHERE id = $2", *req.IsAdmin, id)
	}
	
	return SuccessResponse(c, fiber.Map{"message": "User updated"})
}

// Profile list for admin
type AdminProfile struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	Slug       string    `json:"slug"`
	Name       string    `json:"name"`
	IsActive   bool      `json:"is_active"`
	LinkCount  int       `json:"link_count"`
	TotalClicks int      `json:"total_clicks"`
	CreatedAt  time.Time `json:"created_at"`
}

// ListProfiles returns all profiles for admin moderation
func (h *AdminHandler) ListProfiles(c *fiber.Ctx) error {
	ctx := context.Background()
	
	search := c.Query("search", "")
	
	query := `
		SELECT p.id, p.user_id, u.username, p.slug, p.name, p.is_active,
			COALESCE((SELECT COUNT(*) FROM links WHERE profile_id = p.id), 0) as link_count,
			COALESCE((SELECT SUM(clicks) FROM links WHERE profile_id = p.id), 0) as total_clicks,
			p.created_at
		FROM profiles p
		JOIN users u ON u.id = p.user_id
		WHERE ($1 = '' OR p.name ILIKE '%' || $1 || '%' OR p.slug ILIKE '%' || $1 || '%' OR u.username ILIKE '%' || $1 || '%')
		ORDER BY p.created_at DESC
		LIMIT 100
	`
	
	rows, err := h.db.Query(ctx, query, search)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	defer rows.Close()
	
	profiles := []AdminProfile{}
	for rows.Next() {
		var p AdminProfile
		err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Slug, &p.Name, &p.IsActive, &p.LinkCount, &p.TotalClicks, &p.CreatedAt)
		if err != nil {
			continue
		}
		profiles = append(profiles, p)
	}
	
	return SuccessResponse(c, profiles)
}

// UpdateProfile updates profile status (hide/show)
func (h *AdminHandler) UpdateProfile(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid profile ID")
	}
	
	var req struct {
		IsActive *bool `json:"is_active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}
	
	ctx := context.Background()
	
	if req.IsActive != nil {
		h.db.Exec(ctx, "UPDATE profiles SET is_active = $1 WHERE id = $2", *req.IsActive, id)
	}
	
	return SuccessResponse(c, fiber.Map{"message": "Profile updated"})
}

// GetUserDetail gets detailed info about a user
func (h *AdminHandler) GetUserDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return ValidationError(c, "Invalid user ID")
	}
	
	ctx := context.Background()
	
	// Get user
	var user struct {
		models.User
		ProfileCount int `json:"profile_count"`
	}
	err = h.db.QueryRow(ctx, `
		SELECT id, username, email, is_verified, is_active, is_admin, created_at,
			(SELECT COUNT(*) FROM profiles WHERE user_id = $1) as profile_count
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.IsVerified, &user.IsActive, &user.IsAdmin, &user.CreatedAt, &user.ProfileCount)
	
	if err != nil {
		return NotFound(c, "User")
	}
	
	// Get profiles
	rows, _ := h.db.Query(ctx, `
		SELECT id, slug, name, is_active, 
			(SELECT COUNT(*) FROM links WHERE profile_id = profiles.id) as link_count,
			(SELECT COALESCE(SUM(clicks), 0) FROM links WHERE profile_id = profiles.id) as total_clicks
		FROM profiles WHERE user_id = $1
	`, id)
	defer rows.Close()
	
	profiles := []fiber.Map{}
	for rows.Next() {
		var p struct {
			ID          int    `json:"id"`
			Slug        string `json:"slug"`
			Name        string `json:"name"`
			IsActive    bool   `json:"is_active"`
			LinkCount   int    `json:"link_count"`
			TotalClicks int    `json:"total_clicks"`
		}
		rows.Scan(&p.ID, &p.Slug, &p.Name, &p.IsActive, &p.LinkCount, &p.TotalClicks)
		profiles = append(profiles, fiber.Map{
			"id": p.ID, "slug": p.Slug, "name": p.Name, 
			"is_active": p.IsActive, "link_count": p.LinkCount, "total_clicks": p.TotalClicks,
		})
	}
	
	return SuccessResponse(c, fiber.Map{
		"user":     user.ToPublic(),
		"profiles": profiles,
	})
}
