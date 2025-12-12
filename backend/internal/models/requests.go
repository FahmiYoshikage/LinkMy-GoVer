package models

// RegisterRequest for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse returned after successful auth
type AuthResponse struct {
	User         UserPublic `json:"user"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
}

// RefreshRequest for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// CreateProfileRequest for creating a new profile
type CreateProfileRequest struct {
	Slug  string  `json:"slug" validate:"required,min=3,max=50,alphanum"`
	Name  string  `json:"name" validate:"required,max=100"`
	Title *string `json:"title,omitempty"`
	Bio   *string `json:"bio,omitempty"`
}

// UpdateProfileRequest for updating a profile
type UpdateProfileRequest struct {
	Slug     *string `json:"slug,omitempty"`
	Name     *string `json:"name,omitempty"`
	Title    *string `json:"title,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// CreateLinkRequest for creating a new link
type CreateLinkRequest struct {
	Title      string `json:"title" validate:"required,max=100"`
	URL        string `json:"url" validate:"required,url,max=500"`
	Icon       string `json:"icon,omitempty"`
	CategoryID *int   `json:"category_id,omitempty"`
	Position   *int   `json:"position,omitempty"`
}

// UpdateLinkRequest for updating a link
type UpdateLinkRequest struct {
	Title      *string `json:"title,omitempty"`
	URL        *string `json:"url,omitempty"`
	Icon       *string `json:"icon,omitempty"`
	CategoryID *int    `json:"category_id,omitempty"`
	Position   *int    `json:"position,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// ReorderLinksRequest for reordering links
type ReorderLinksRequest struct {
	Links []LinkPosition `json:"links" validate:"required"`
}

// LinkPosition for reordering
type LinkPosition struct {
	ID       int `json:"id"`
	Position int `json:"position"`
}

// CreateCategoryRequest for creating a category
type CreateCategoryRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Icon     string `json:"icon,omitempty"`
	Color    string `json:"color,omitempty"`
	Position *int   `json:"position,omitempty"`
}

// UpdateCategoryRequest for updating a category
type UpdateCategoryRequest struct {
	Name       *string `json:"name,omitempty"`
	Icon       *string `json:"icon,omitempty"`
	Color      *string `json:"color,omitempty"`
	Position   *int    `json:"position,omitempty"`
	IsExpanded *bool   `json:"is_expanded,omitempty"`
}

// UpdateThemeRequest for updating theme
type UpdateThemeRequest struct {
	BgType             *string `json:"bg_type,omitempty"`
	BgValue            *string `json:"bg_value,omitempty"`
	ButtonStyle        *string `json:"button_style,omitempty"`
	ButtonColor        *string `json:"button_color,omitempty"`
	TextColor          *string `json:"text_color,omitempty"`
	Font               *string `json:"font,omitempty"`
	Layout             *string `json:"layout,omitempty"`
	ContainerStyle     *string `json:"container_style,omitempty"`
	EnableAnimations   *bool   `json:"enable_animations,omitempty"`
	EnableGlassEffect  *bool   `json:"enable_glass_effect,omitempty"`
	ShadowIntensity    *string `json:"shadow_intensity,omitempty"`
	BoxedEnabled       *bool   `json:"boxed_enabled,omitempty"`
	BoxedOuterBgType   *string `json:"boxed_outer_bg_type,omitempty"`
	BoxedOuterBgValue  *string `json:"boxed_outer_bg_value,omitempty"`
	BoxedContainerBg   *string `json:"boxed_container_bg,omitempty"`
	BoxedMaxWidth      *int    `json:"boxed_max_width,omitempty"`
	BoxedRadius        *int    `json:"boxed_radius,omitempty"`
	BoxedShadow        *bool   `json:"boxed_shadow,omitempty"`
}

// TrackClickRequest for tracking link clicks
type TrackClickRequest struct {
	Referrer *string `json:"referrer,omitempty"`
}

// PublicProfile is the response for public profile viewing
type PublicProfile struct {
	Profile    Profile    `json:"profile"`
	Theme      Theme      `json:"theme"`
	Categories []Category `json:"categories"`
	Links      []Link     `json:"links"`
	IsVerified bool       `json:"is_verified"`
}

// AnalyticsResponse for profile analytics
type AnalyticsResponse struct {
	TotalClicks     int                   `json:"total_clicks"`
	ClicksByDay     []DayStats            `json:"clicks_by_day"`
	ClicksByLink    []LinkStats           `json:"clicks_by_link"`
	ClicksByCountry []CountryStats        `json:"clicks_by_country"`
	TopReferrers    []ReferrerStats       `json:"top_referrers"`
}

// DayStats for daily click stats
type DayStats struct {
	Date   string `json:"date"`
	Clicks int    `json:"clicks"`
}

// LinkStats for per-link stats
type LinkStats struct {
	LinkID int    `json:"link_id"`
	Title  string `json:"title"`
	Clicks int    `json:"clicks"`
}

// CountryStats for geographic stats
type CountryStats struct {
	Country string `json:"country"`
	Clicks  int    `json:"clicks"`
}

// ReferrerStats for referrer tracking
type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Clicks   int    `json:"clicks"`
}
