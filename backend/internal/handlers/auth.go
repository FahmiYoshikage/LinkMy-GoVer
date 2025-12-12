package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/FahmiYoshikage/linkmy-v2/internal/config"
	"github.com/FahmiYoshikage/linkmy-v2/internal/middleware"
	"github.com/FahmiYoshikage/linkmy-v2/internal/models"
	"github.com/FahmiYoshikage/linkmy-v2/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	profileRepo *repository.ProfileRepository
	cfg         *config.Config
}

func NewAuthHandler(db *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo:    repository.NewUserRepository(db),
		sessionRepo: repository.NewSessionRepository(db),
		profileRepo: repository.NewProfileRepository(db),
		cfg:         cfg,
	}
}

// SendOTP sends OTP to email for registration
func (h *AuthHandler) SendOTP(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	if len(req.Password) < 8 {
		return ValidationError(c, "Password must be at least 8 characters")
	}

	ctx := context.Background()

	// Check if email exists
	exists, err := h.userRepo.ExistsEmail(ctx, req.Email)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if exists {
		return ErrorResponse(c, fiber.StatusConflict, "Email already registered")
	}

	// Generate and store OTP
	otp := GenerateOTP()
	StoreOTP(req.Email, otp)

	// Send email
	if err := SendOTPEmail(req.Email, otp); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send OTP email")
	}

	return SuccessResponse(c, fiber.Map{
		"message": "OTP sent to email",
		"email":   req.Email,
	})
}

// VerifyOTPEndpoint verifies OTP code
func (h *AuthHandler) VerifyOTPEndpoint(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	if !VerifyOTP(req.Email, req.OTP) {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired OTP")
	}

	return SuccessResponse(c, fiber.Map{
		"message":  "OTP verified",
		"verified": true,
	})
}

// CompleteRegistration finishes registration after OTP verification
func (h *AuthHandler) CompleteRegistration(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OTP      string `json:"otp"`
		Username string `json:"username"`
		Slug     string `json:"slug"`
	}
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	// Verify OTP first
	if !VerifyOTP(req.Email, req.OTP) {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired OTP")
	}

	// Validation
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return ValidationError(c, "Username must be 3-50 characters")
	}
	if len(req.Password) < 8 {
		return ValidationError(c, "Password must be at least 8 characters")
	}

	ctx := context.Background()

	// Check if username exists
	exists, err := h.userRepo.ExistsUsername(ctx, req.Username)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if exists {
		return ErrorResponse(c, fiber.StatusConflict, "Username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsVerified:   true, // Verified via OTP
		IsActive:     true,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrorResponse(c, fiber.StatusConflict, "User already exists")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	// Create profile with slug
	slug := req.Slug
	if slug == "" {
		slug = req.Username
	}
	profile := &models.Profile{
		UserID:   user.ID,
		Slug:     slug,
		Name:     req.Username,
		Title:    &req.Username,
		Avatar:   "default-avatar.png",
		IsActive: true,
	}
	if err := h.profileRepo.Create(ctx, profile); err != nil {
		// Log error but don't fail
	}

	// Delete OTP after successful registration
	DeleteOTP(req.Email)

	// Generate tokens
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	refreshToken, err := h.createSession(ctx, user, c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		User:         user.ToPublic(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.cfg.JWTExpiryHours * 3600),
	})
}

// Register creates a new user account (legacy, still works)
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	// Basic validation
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return ValidationError(c, "Username must be 3-50 characters")
	}
	if len(req.Password) < 8 {
		return ValidationError(c, "Password must be at least 8 characters")
	}

	ctx := context.Background()

	// Check if email exists
	exists, err := h.userRepo.ExistsEmail(ctx, req.Email)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if exists {
		return ErrorResponse(c, fiber.StatusConflict, "Email already registered")
	}

	// Check if username exists
	exists, err = h.userRepo.ExistsUsername(ctx, req.Username)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}
	if exists {
		return ErrorResponse(c, fiber.StatusConflict, "Username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsVerified:   false, // Require email verification
		IsActive:     true,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrorResponse(c, fiber.StatusConflict, "User already exists")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user")
	}

	// Create default profile with username as slug
	profile := &models.Profile{
		UserID:   user.ID,
		Slug:     req.Username,
		Name:     req.Username + " - Main Profile",
		Title:    &req.Username,
		Avatar:   "default-avatar.png",
		IsActive: true,
	}
	if err := h.profileRepo.Create(ctx, profile); err != nil {
		// Log error but don't fail registration
		// Profile can be created later
	}

	// Generate tokens
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	refreshToken, err := h.createSession(ctx, user, c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		User:         user.ToPublic(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.cfg.JWTExpiryHours * 3600),
	})
}

// Login authenticates a user
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	ctx := context.Background()

	// Find user by email
	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Database error")
	}

	// Check if user is active
	if !user.IsActive {
		return ErrorResponse(c, fiber.StatusForbidden, "Account is disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password")
	}

	// Generate tokens
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	refreshToken, err := h.createSession(ctx, user, c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create session")
	}

	return SuccessResponse(c, models.AuthResponse{
		User:         user.ToPublic(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.cfg.JWTExpiryHours * 3600),
	})
}

// RefreshToken generates new access token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	ctx := context.Background()

	// Find session
	session, err := h.sessionRepo.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token")
	}

	// Get user
	user, err := h.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	// Generate new access token
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	return SuccessResponse(c, fiber.Map{
		"access_token": accessToken,
		"expires_in":   int64(h.cfg.JWTExpiryHours * 3600),
	})
}

// Logout invalidates the user's session
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req models.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		// Even if no body, we can still logout
		return SuccessResponse(c, fiber.Map{"message": "Logged out"})
	}

	ctx := context.Background()

	// Find and delete session
	session, err := h.sessionRepo.GetByRefreshToken(ctx, req.RefreshToken)
	if err == nil && session != nil {
		h.sessionRepo.Delete(ctx, session.ID)
	}

	return SuccessResponse(c, fiber.Map{"message": "Logged out successfully"})
}

// GetCurrentUser returns the current user's info
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return Unauthorized(c)
	}

	ctx := context.Background()
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return NotFound(c, "User")
	}

	return SuccessResponse(c, user.ToPublic())
}

// UpdateCurrentUser updates the current user's info
func (h *AuthHandler) UpdateCurrentUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return Unauthorized(c)
	}

	var req struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return ValidationError(c, "Invalid request body")
	}

	ctx := context.Background()
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return NotFound(c, "User")
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := h.userRepo.Update(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return ErrorResponse(c, fiber.StatusConflict, "Username or email already taken")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user")
	}

	return SuccessResponse(c, user.ToPublic())
}

// Helper: Generate access token
func (h *AuthHandler) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Duration(h.cfg.JWTExpiryHours) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

// Helper: Create session with refresh token
func (h *AuthHandler) createSession(ctx context.Context, user *models.User, c *fiber.Ctx) (string, error) {
	// Generate refresh token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	refreshToken := hex.EncodeToString(tokenBytes)

	ip := c.IP()
	userAgent := c.Get("User-Agent")

	session := &models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IP:           &ip,
		UserAgent:    &userAgent,
		ExpiresAt:    time.Now().Add(time.Duration(h.cfg.RefreshExpiryHours) * time.Hour),
	}

	if err := h.sessionRepo.Create(ctx, session); err != nil {
		return "", err
	}

	return refreshToken, nil
}
