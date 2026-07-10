package main

import (
	"context"
	"os"
	"social-engine/common/logger"
	"social-engine/common/repositories"
	"social-engine/common/validation"
	"social-engine/handlers"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to the database
	err := repositories.Connect(ctx)
	if err != nil {
		logger.L(ctx).Fatal("Failed to connect to repositories", zap.Error(err))
	}

	// Initialize the validator
	validation.InitValidator()

	// Set up CORS allowed origins
	allowedOrigins := make(map[string]struct{})
	for o := range strings.SplitSeq(os.Getenv("CORS_ALLOWED_ORIGINS"), ",") {
		if o = strings.ToLower(strings.TrimSpace(o)); o != "" {
			allowedOrigins[o] = struct{}{}
		}
	}

	// Initialize Fiber
	app := fiber.New()
	app.Use(
		cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				_, ok := allowedOrigins[origin]
				return ok
			},
			AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:  "Origin,Content-Type,Accept,Authorization",
			ExposeHeaders: "Content-Length",
			MaxAge:        3600,
		}),
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Social Engine API!")
	})

	app.Get("/docs/*", swagger.New(swagger.ConfigDefault))

	app.Get("/health", handlers.Health)

	/* Initiate server */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	zap.L().Fatal("Unable to listen", zap.Error(app.Listen(":"+port)))
}
