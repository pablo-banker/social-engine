package middleware

import (
	"context"
	"strings"

	"social-engine/common/apiErrors"
	"social-engine/common/models"
	"social-engine/common/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// bearerToken extracts the token from an "Authorization: Bearer <token>" header.
func bearerToken(c *fiber.Ctx) string {
	header := c.Get(fiber.HeaderAuthorization)
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// setUser stores the authenticated user ID on the request Locals (for handlers)
// and on the UserContext (so the logger can attach it to every log line).
func setUser(c *fiber.Ctx, userID uuid.UUID) {
	c.Locals("userId", userID)
	c.SetUserContext(context.WithValue(c.UserContext(), "user_id", userID.String()))
}

// RequireAuth rejects requests that lack a valid Bearer token.
func RequireAuth(c *fiber.Ctx) error {
	token := bearerToken(c)
	if token == "" {
		return utils.BuildErrorResponse(c, apiErrors.ErrUnauthorized)
	}

	userID, err := models.ParseToken(token)
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidToken.WithInternal(err.Error()))
	}

	setUser(c, userID)
	return c.Next()
}

// OptionalAuth attaches the user when a valid token is present but lets
// anonymous requests through. Used by public reads (feed, profile) that
// compute per-user state such as likedByMe when the caller is logged in.
func OptionalAuth(c *fiber.Ctx) error {
	if token := bearerToken(c); token != "" {
		if userID, err := models.ParseToken(token); err == nil {
			setUser(c, userID)
		}
	}
	return c.Next()
}
