package utils

import "github.com/gofiber/fiber/v2"

func Response(c *fiber.Ctx, data any, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"data":    data,
		"status":  status,
		"message": message,
	})
}
