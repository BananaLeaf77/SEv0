package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	Messages []string `json:"messages"`
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
			defaultRepeater := 1
			payload.Repeater = &defaultRepeater
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		done := make(chan error, 1)

		go func() {
			// Loop through recipients
			for _, phone := range payload.To {
				completeFormat := fmt.Sprintf("%s%s", phone[:2], phone[2:])
				JID := types.NewJID(completeFormat, types.DefaultUserServer)

				// Loop through messages
				for _, msg := range payload.Messages {
					ConversationMessage := &waE2E.Message{
						Conversation: &msg,
					}

					// Repeat message N times
					for i := 0; i < *payload.Repeater; i++ {
						_, err := WA.SendMessage(timeoutCtx, JID, ConversationMessage)
						if err != nil {
							done <- fmt.Errorf("failed to send message to %s: %w", phone, err)
							return
						}
					}
				}
			}
			done <- nil
		}()

		select {
		case <-timeoutCtx.Done():
			log.Print("/send route timed out, code: ", utils.ColorStatus(500))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Request timed out while sending messages, you should make sure the phone numbers are correct :)",
				"code":  500,
			})
		case err := <-done:
			if err != nil {
				log.Print("/send route accessed failed, code: ", utils.ColorStatus(500), ", error: ", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
					"code":  500,
				})
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
