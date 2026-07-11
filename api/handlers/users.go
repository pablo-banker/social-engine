package handlers

import (
	"errors"
	"net/http"
	"strings"

	"social-engine/common/apiErrors"
	"social-engine/common/models"
	"social-engine/common/repositories"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetProfile returns a user's public profile with aggregate stats.
// @Summary Get a user profile
// @Tags Users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} entities.User
// @Failure 404 {object} utils.ErrorResponse
// @Router /users/{username} [get]
func GetProfile(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	username := strings.ToLower(strings.TrimSpace(c.Params("username")))

	var user entities.User
	err := repo.FindOne(ctx, &user, &entities.QueryParams{
		Query: entities.Query{Filters: "lower(username) = ?", Values: []any{username}},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.BuildErrorResponse(c, apiErrors.ErrUserNotFound)
		}
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	posts, err := repo.Count(ctx, &entities.Post{}, &entities.QueryParams{
		Query: entities.Query{Filters: "author_id = ?", Values: []any{user.ID}},
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	followers, err := repo.Count(ctx, &entities.Follow{}, &entities.QueryParams{
		Query: entities.Query{Filters: "following_id = ?", Values: []any{user.ID}},
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	following, err := repo.Count(ctx, &entities.Follow{}, &entities.QueryParams{
		Query: entities.Query{Filters: "follower_id = ?", Values: []any{user.ID}},
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	user.Stats = &entities.UserStats{
		Posts:     int(posts),
		Followers: int(followers),
		Following: int(following),
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, &user)
}

// ListUserPosts returns a user's posts, newest first.
// @Summary List a user's posts
// @Tags Users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {array} entities.Post
// @Router /users/{username}/posts [get]
func ListUserPosts(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	username := strings.ToLower(strings.TrimSpace(c.Params("username")))

	var user entities.User
	err := repo.FindOne(ctx, &user, &entities.QueryParams{
		Query: entities.Query{Filters: "lower(username) = ?", Values: []any{username}},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.BuildSuccessResponse(c, http.StatusOK, []*entities.Post{})
		}
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	result, err := repo.FindAll(ctx, &entities.Post{}, &entities.QueryParams{
		Query: entities.Query{Filters: "author_id = ?", Values: []any{user.ID}},
		Sort:  "created_at desc",
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, models.DecoratePosts(result, viewer(c)))
}
