package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth middleware for protected routes
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Missing authorization header",
			})
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Invalid token claims",
			})
		}

		// Store user ID in context
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "unauthorized",
				"message": "Invalid user ID in token",
			})
		}

		c.Locals("userID", int(userID))
		c.Locals("username", claims["username"])
		
		// Store is_admin if present
		if isAdmin, ok := claims["is_admin"].(bool); ok {
			c.Locals("is_admin", isAdmin)
		} else {
			c.Locals("is_admin", false)
		}

		return c.Next()
	}
}

// GetUserID extracts user ID from context (use after JWTAuth middleware)
func GetUserID(c *fiber.Ctx) int {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return 0
	}
	return userID
}

// GetUsername extracts username from context
func GetUsername(c *fiber.Ctx) string {
	username, ok := c.Locals("username").(string)
	if !ok {
		return ""
	}
	return username
}
