package handlers

import (
	"context"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/middleware"
	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/FahmiYoshikage/linkmy-v2/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsHandler struct {
	db          *pgxpool.Pool
	profileRepo *repository.ProfileRepository
}

func NewAnalyticsHandler(db *pgxpool.Pool) *AnalyticsHandler {
	return &AnalyticsHandler{
		db:          db,
		profileRepo: repository.NewProfileRepository(db),
	}
}

// GetProfileAnalytics returns analytics for a profile
func (h *AnalyticsHandler) GetProfileAnalytics(c *fiber.Ctx) error {
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

	// Get time range from query params (default: last 30 days)
	days := c.QueryInt("days", 30)
	startDate := time.Now().AddDate(0, 0, -days)

	analytics := &models.AnalyticsResponse{}

	// Total clicks
	var totalClicks int
	h.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(l.clicks), 0)
		FROM links l WHERE l.profile_id = $1
	`, profileID).Scan(&totalClicks)
	analytics.TotalClicks = totalClicks

	// Clicks by day
	rows, err := h.db.Query(ctx, `
		SELECT DATE(c.clicked_at) as date, COUNT(*) as clicks
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.profile_id = $1 AND c.clicked_at >= $2
		GROUP BY DATE(c.clicked_at)
		ORDER BY date ASC
	`, profileID, startDate)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ds models.DayStats
			var date time.Time
			rows.Scan(&date, &ds.Clicks)
			ds.Date = date.Format("2006-01-02")
			analytics.ClicksByDay = append(analytics.ClicksByDay, ds)
		}
	}

	// Clicks by link
	rows, err = h.db.Query(ctx, `
		SELECT l.id, l.title, l.clicks
		FROM links l
		WHERE l.profile_id = $1 AND l.clicks > 0
		ORDER BY l.clicks DESC
		LIMIT 10
	`, profileID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ls models.LinkStats
			rows.Scan(&ls.LinkID, &ls.Title, &ls.Clicks)
			analytics.ClicksByLink = append(analytics.ClicksByLink, ls)
		}
	}

	// Clicks by country
	rows, err = h.db.Query(ctx, `
		SELECT COALESCE(c.country, 'Unknown') as country, COUNT(*) as clicks
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.profile_id = $1 AND c.clicked_at >= $2
		GROUP BY c.country
		ORDER BY clicks DESC
		LIMIT 10
	`, profileID, startDate)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cs models.CountryStats
			rows.Scan(&cs.Country, &cs.Clicks)
			analytics.ClicksByCountry = append(analytics.ClicksByCountry, cs)
		}
	}

	// Top referrers
	rows, err = h.db.Query(ctx, `
		SELECT COALESCE(c.referrer, 'Direct') as referrer, COUNT(*) as clicks
		FROM clicks c
		JOIN links l ON c.link_id = l.id
		WHERE l.profile_id = $1 AND c.clicked_at >= $2
		GROUP BY c.referrer
		ORDER BY clicks DESC
		LIMIT 10
	`, profileID, startDate)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rs models.ReferrerStats
			rows.Scan(&rs.Referrer, &rs.Clicks)
			analytics.TopReferrers = append(analytics.TopReferrers, rs)
		}
	}

	return SuccessResponse(c, analytics)
}
