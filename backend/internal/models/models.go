package models

import "time"

// User represents a registered user
type User struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never expose
	IsVerified   bool       `json:"is_verified"`
	IsActive     bool       `json:"is_active"`
	IsAdmin      bool       `json:"is_admin"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// UserPublic is the safe version for API responses
type UserPublic struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToPublic converts User to UserPublic
func (u *User) ToPublic() UserPublic {
	return UserPublic{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		IsAdmin:    u.IsAdmin,
		CreatedAt:  u.CreatedAt,
	}
}

// Profile represents a user's profile page
type Profile struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	Slug         string     `json:"slug"`
	Name         string     `json:"name"`
	Title        *string    `json:"title,omitempty"`
	Bio          *string    `json:"bio,omitempty"`
	Avatar       string     `json:"avatar"`
	IsActive     bool       `json:"is_active"`
	DisplayOrder int        `json:"display_order"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// ProfileWithStats includes link count and click stats
type ProfileWithStats struct {
	Profile
	LinkCount   int `json:"link_count"`
	TotalClicks int `json:"total_clicks"`
}

// Theme represents profile appearance settings
type Theme struct {
	ID                 int        `json:"id"`
	ProfileID          int        `json:"profile_id"`
	BgType             string     `json:"bg_type"`
	BgValue            *string    `json:"bg_value,omitempty"`
	ButtonStyle        string     `json:"button_style"`
	ButtonColor        string     `json:"button_color"`
	TextColor          string     `json:"text_color"`
	Font               string     `json:"font"`
	Layout             string     `json:"layout"`
	ContainerStyle     string     `json:"container_style"`
	EnableAnimations   bool       `json:"enable_animations"`
	EnableGlassEffect  bool       `json:"enable_glass_effect"`
	ShadowIntensity    string     `json:"shadow_intensity"`
	BoxedEnabled       bool       `json:"boxed_enabled"`
	BoxedOuterBgType   *string    `json:"boxed_outer_bg_type,omitempty"`
	BoxedOuterBgValue  *string    `json:"boxed_outer_bg_value,omitempty"`
	BoxedContainerBg   string     `json:"boxed_container_bg"`
	BoxedMaxWidth      int        `json:"boxed_max_width"`
	BoxedRadius        int        `json:"boxed_radius"`
	BoxedShadow        bool       `json:"boxed_shadow"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

// Category represents a link category/folder
type Category struct {
	ID         int       `json:"id"`
	ProfileID  int       `json:"profile_id"`
	Name       string    `json:"name"`
	Icon       string    `json:"icon"`
	Color      string    `json:"color"`
	Position   int       `json:"position"`
	IsExpanded bool      `json:"is_expanded"`
	CreatedAt  time.Time `json:"created_at"`
}

// Link represents a profile link
type Link struct {
	ID         int        `json:"id"`
	ProfileID  int        `json:"profile_id"`
	CategoryID *int       `json:"category_id,omitempty"`
	Title      string     `json:"title"`
	URL        string     `json:"url"`
	Icon       string     `json:"icon"`
	Position   int        `json:"position"`
	Clicks     int        `json:"clicks"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// Click represents a link click event
type Click struct {
	ID        int64     `json:"id"`
	LinkID    int       `json:"link_id"`
	IP        *string   `json:"ip,omitempty"`
	Country   *string   `json:"country,omitempty"`
	City      *string   `json:"city,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
	Referrer  *string   `json:"referrer,omitempty"`
	ClickedAt time.Time `json:"clicked_at"`
}

// Session represents a user session
type Session struct {
	ID           string    `json:"id"`
	UserID       int       `json:"user_id"`
	RefreshToken string    `json:"-"`
	IP           *string   `json:"ip,omitempty"`
	UserAgent    *string   `json:"user_agent,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// SocialIcon represents a social media platform
type SocialIcon struct {
	ID           int     `json:"id"`
	PlatformName string  `json:"platform_name"`
	IconClass    string  `json:"icon_class"`
	IconColor    *string `json:"icon_color,omitempty"`
	BaseURL      *string `json:"base_url,omitempty"`
}
