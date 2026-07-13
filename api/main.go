package main

import (
	"context"
	"os"
	"social-engine/common/apiErrors"
	"social-engine/common/logger"
	"social-engine/common/models"
	"social-engine/common/repositories"
	"social-engine/common/utils"
	"social-engine/common/validation"
	"social-engine/handlers"
	"social-engine/handlers/middleware"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to the database
	err := repositories.Connect(ctx)
	if err != nil {
		logger.L(ctx).Fatal("Failed to connect to repositories", zap.Error(err))
	}

	// Fail fast if the JWT signing secret is missing or too weak.
	if err := models.ValidateSecret(); err != nil {
		logger.L(ctx).Fatal("Invalid JWT configuration", zap.Error(err))
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

	// Throttle the auth endpoints per client IP to blunt brute-force and
	// credential-stuffing attempts.
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return utils.BuildErrorResponse(c, apiErrors.ErrTooManyRequests)
		},
	})

	// Auth
	app.Post("/auth/register", authLimiter, handlers.Register)
	app.Post("/auth/login", authLimiter, handlers.Login)

	// Posts
	app.Get("/posts", middleware.OptionalAuth, handlers.List)
	app.Post("/posts", middleware.RequireAuth, handlers.Create)
	app.Get("/posts/:id", middleware.OptionalAuth, handlers.Get)
	app.Post("/posts/:id/like", middleware.RequireAuth, handlers.ToggleLike)
	app.Get("/posts/:id/comments", handlers.ListComments)
	app.Post("/posts/:id/comments", middleware.RequireAuth, handlers.AddComment)

	// Users / Profiles
	app.Get("/users/:username", handlers.GetProfile)
	app.Get("/users/:username/posts", middleware.OptionalAuth, handlers.ListUserPosts)

	// Explore & Trending
	app.Get("/explore", middleware.OptionalAuth, handlers.Explore)
	app.Get("/trending", middleware.OptionalAuth, handlers.Trending)

	// Me (logged-in user)
	app.Get("/me", middleware.RequireAuth, handlers.GetSettings)
	app.Patch("/me", middleware.RequireAuth, handlers.UpdateProfile)

	/* Initiate server */
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	zap.L().Fatal("Unable to listen", zap.Error(app.Listen(":"+port)))
}
