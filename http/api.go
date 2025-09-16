package api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"

	"SEv0/middleware"
	"SEv0/utils"
)

type Payload struct {
	To       []string `json:"to"`
	Message  string   `json:"message"`
	Repeater *int     `json:"repeater"`
}

func InitApi(WA *whatsmeow.Client, ctx context.Context) error {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOW_ORIGINS"),
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/ping", func(c *fiber.Ctx) error {
		log.Print("/ping route accessed successfully, code: ", utils.ColorStatus(200))
		return c.JSON(fiber.Map{
			"message": "Backend Served",
			"code":    200,
		})

	})

	app.Post("/send", middleware.ValidatePayloadIdentity(), func(c *fiber.Ctx) error {
		var payload Payload
		if err := c.BodyParser(&payload); err != nil {
			log.Print("/send route accessed failed, code: ", utils.ColorStatus(400), " error: ", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
				"code":  400,
			})
		}

		if payload.Repeater == nil {
			*payload.Repeater = 1
		}

		for _, phone := range payload.To {
			completeFormat := fmt.Sprintf("%s%s", phone[:2], phone[2:])
			JID := types.NewJID(completeFormat, types.DefaultUserServer)

			ConversationMessage := &waE2E.Message{
				Conversation: &payload.Message,
			}

			for i := 0; i < *payload.Repeater; i++ {
				_, err := WA.SendMessage(ctx, JID, ConversationMessage)
				if err != nil {
					log.Print("/send route accessed failed, code: ", utils.ColorStatus(500), ", error: ", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to send message to " + phone,
						"code":  500,
					})
				}
			}
		}

		log.Print("/send route accessed successfully, code: ", utils.ColorStatus(200))
		return c.JSON(fiber.Map{
			"status": "Message(s) sent",
			"code":   200,
		})
	})

	return app.Listen(":3000")
}
