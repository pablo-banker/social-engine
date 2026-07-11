package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// fixedTime keeps decorated-row timestamps deterministic across tests.
var fixedTime = time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

// withUser is a test middleware that authenticates the request as the given
// user, mirroring what middleware.RequireAuth does at runtime.
func withUser(id uuid.UUID) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("userId", id)
		return c.Next()
	}
}
