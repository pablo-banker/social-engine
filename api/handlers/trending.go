package handlers

import (
	"net/http"

	"social-engine/common/apiErrors"
	"social-engine/common/models"
	"social-engine/common/repositories"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/gofiber/fiber/v2"
)

type trendingResponse struct {
	Topics []models.TrendingTopic `json:"topics"`
	Posts  []*entities.Post       `json:"posts"`
}

// Trending returns the hottest hashtags and the most-liked posts.
// @Summary Trending topics and posts
// @Tags Trending
// @Produce json
// @Success 200 {object} trendingResponse
// @Router /trending [get]
func Trending(c *fiber.Ctx) error {
	var (
		ctx      = c.UserContext()
		repo     = repositories.InitializeRepoInstance(ctx)
		viewerID = viewer(c)
	)

	topics, err := models.TrendingTopics(ctx, repo, trendingTopicsLimit)
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	postsResult, err := repo.FindAll(ctx, &entities.Post{}, &entities.QueryParams{Sort: "created_at desc"})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}
	posts := models.TopByLikes(models.DecoratePosts(postsResult, viewerID), popularPostsLimit)

	return utils.BuildSuccessResponse(c, http.StatusOK, trendingResponse{
		Topics: topics,
		Posts:  posts,
	})
}
