package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Limits for list endpoints.
const (
	suggestedUsersLimit = 6
	popularPostsLimit   = 10
	trendingTopicsLimit = 8
)

// viewer returns the authenticated user's id, or uuid.Nil for anonymous
// callers. It is populated by the auth middleware via c.Locals("userId").
func viewer(c *fiber.Ctx) uuid.UUID {
	id, _ := c.Locals("userId").(uuid.UUID)
	return id
}
