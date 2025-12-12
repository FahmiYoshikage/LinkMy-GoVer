package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AdminAuth middleware checks if user is an admin
func AdminAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin := c.Locals("is_admin")
		
		if isAdmin == nil || !isAdmin.(bool) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Admin access required",
			})
		}
		
		return c.Next()
	}
}
