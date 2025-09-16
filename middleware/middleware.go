package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func ValidatePayloadIdentity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token != "Bearer your_secret_token" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		return c.Next()
	}
}
