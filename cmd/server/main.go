package main

import (
	"fmt"
	"game-engine/internal/engine"
	"game-engine/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize Game Engine
	ge := engine.NewGameEngine()

	// Metrics Reporter: Print stats every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			total, correct, incorrect, discoveryTime := ge.GetMetrics()
			fmt.Printf("\n--- Metrics Report (Every 5s) ---\n")
			fmt.Printf("Total Submitted: %d\n", total)
			fmt.Printf("Total Correct:   %d\n", correct)
			fmt.Printf("Total Incorrect: %d\n", incorrect)

			if discoveryTime > 0 {
				fmt.Printf("Discovery Time:  %v\n", discoveryTime)
			} else {
				fmt.Printf("Discovery Time:  N/A (No winner yet)\n")
			}
			fmt.Printf("---------------------------------\n")
		}
	}()

	// Initialize Fiber app
	app := fiber.New()

	// API Endpoint for user responses
	app.Post("/submit", func(c *fiber.Ctx) error {
		var resp models.UserResponse

		// Parse and validate JSON body
		if err := c.BodyParser(&resp); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON body",
			})
		}

		// Validation: Ensure UserID is not empty
		if resp.UserID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "user_id is required",
			})
		}

		if resp.Answer != models.AnswerYes && resp.Answer != models.AnswerNo {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "answer must be yes or no",
			})
		}

		// Forward to Game Engine
		_ = ge.Evaluate(resp)

		return c.SendStatus(fiber.StatusAccepted)
	})

	port := 8080
	fmt.Printf("Fiber API Server starting on :%d...\n", port)
	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
