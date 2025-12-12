package main

import (
	"log"
	"os"

	"github.com/FahmiYoshikage/linkmy-v2/internal/config"
	"github.com/FahmiYoshikage/linkmy-v2/internal/database"
	"github.com/FahmiYoshikage/linkmy-v2/internal/handlers"
	"github.com/FahmiYoshikage/linkmy-v2/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "LinkMy API v2.0",
		ServerHeader: "LinkMy",
		ErrorHandler: handlers.ErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "2.0.0",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Public routes
	authHandler := handlers.NewAuthHandler(db, cfg)
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/refresh", authHandler.RefreshToken)
	api.Post("/auth/logout", authHandler.Logout)
	
	// OTP routes for multi-step registration
	api.Post("/auth/send-otp", authHandler.SendOTP)
	api.Post("/auth/verify-otp", authHandler.VerifyOTPEndpoint)
	api.Post("/auth/complete-registration", authHandler.CompleteRegistration)

	// Public profile view
	profileHandler := handlers.NewProfileHandler(db)
	api.Get("/p/:slug", profileHandler.GetPublicProfile)

	// Click tracking (public)
	linkHandler := handlers.NewLinkHandler(db)
	api.Post("/click/:id", linkHandler.TrackClick)

	// Protected routes
	protected := api.Group("/", middleware.JWTAuth(cfg.JWTSecret))

	// User routes
	protected.Get("/me", authHandler.GetCurrentUser)
	protected.Put("/me", authHandler.UpdateCurrentUser)

	// Profile management
	protected.Get("/profiles", profileHandler.GetUserProfiles)
	protected.Post("/profiles", profileHandler.CreateProfile)
	protected.Get("/profiles/:id", profileHandler.GetProfile)
	protected.Put("/profiles/:id", profileHandler.UpdateProfile)
	protected.Delete("/profiles/:id", profileHandler.DeleteProfile)

	// Link management
	protected.Get("/profiles/:profileId/links", linkHandler.GetLinks)
	protected.Post("/profiles/:profileId/links", linkHandler.CreateLink)
	protected.Put("/links/:id", linkHandler.UpdateLink)
	protected.Delete("/links/:id", linkHandler.DeleteLink)
	protected.Put("/links/reorder", linkHandler.ReorderLinks)

	// Category management
	categoryHandler := handlers.NewCategoryHandler(db)
	protected.Get("/profiles/:profileId/categories", categoryHandler.GetCategories)
	protected.Post("/profiles/:profileId/categories", categoryHandler.CreateCategory)
	protected.Put("/categories/:id", categoryHandler.UpdateCategory)
	protected.Delete("/categories/:id", categoryHandler.DeleteCategory)

	// Theme management
	themeHandler := handlers.NewThemeHandler(db)
	protected.Get("/profiles/:profileId/theme", themeHandler.GetTheme)
	protected.Put("/profiles/:profileId/theme", themeHandler.UpdateTheme)

	// Analytics
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	protected.Get("/profiles/:profileId/analytics", analyticsHandler.GetProfileAnalytics)

	// Admin routes (requires JWT + admin check)
	adminHandler := handlers.NewAdminHandler(db)
	admin := api.Group("/admin", middleware.JWTAuth(cfg.JWTSecret), middleware.AdminAuth())
	admin.Get("/stats", adminHandler.GetStats)
	admin.Get("/users", adminHandler.ListUsers)
	admin.Get("/users/:id", adminHandler.GetUserDetail)
	admin.Put("/users/:id", adminHandler.UpdateUser)
	admin.Get("/profiles", adminHandler.ListProfiles)
	admin.Put("/profiles/:id", adminHandler.UpdateProfile)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 LinkMy API starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
