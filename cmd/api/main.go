package main

import (
	"log"

	"appointment-scheduling-bot/internal/shared/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add logger middleware
	app.Use(logger.New())

	// Health check endpoint
	app.Get("/healthz", func(c *fiber.Ctx) error {
		envSummary := fiber.Map{
			"status":     "ok",
			"env":        cfg.Env,
			"http_port":  cfg.HttpPort,
			"timezone":   cfg.Timezone,
			"gcal_calendar_id": cfg.GCalCalendarID,
			"has_google_creds": len(cfg.GoogleCredsJSON) > 0,
			"has_supabase_url": len(cfg.SupabaseURL) > 0,
			"has_supabase_key": len(cfg.SupabaseKey) > 0,
			"redis_url":  cfg.RedisURL,
		}
		return c.JSON(envSummary)
	})

	// Start server
	port := ":" + cfg.HttpPort
	log.Printf("Starting server on port %s", port)
	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 