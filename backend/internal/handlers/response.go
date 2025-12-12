package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is the global error handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default status code
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}

// Response helpers
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}

func ValidationError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error":   "validation_error",
		"message": message,
	})
}

func NotFound(c *fiber.Ctx, resource string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":   "not_found",
		"message": resource + " not found",
	})
}

func Unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   "unauthorized",
		"message": "You are not authorized to perform this action",
	})
}

func Forbidden(c *fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error":   "forbidden",
		"message": "You do not have permission to access this resource",
	})
}
