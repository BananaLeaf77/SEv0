package middleware

import (
	"SEv0/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ValidatePayloadIdentity() fiber.Handler {
	secret := os.Getenv("API_SECRET")

	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token != "Bearer "+secret {
			log.Print("/send route accessed failed, code: ", utils.ColorStatus(401), " Unauthorized", " error: Invalid or missing token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
				"code":  401,
			})
		}
		return c.Next()
	}
}
