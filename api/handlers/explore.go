package handlers

import (
	"net/http"
	"strings"

	"social-engine/common/apiErrors"
	"social-engine/common/models"
	"social-engine/common/repositories"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/gofiber/fiber/v2"
)

type exploreResponse struct {
	Users []*entities.User `json:"users"`
	Posts []*entities.Post `json:"posts"`
}

// Explore returns suggested users and popular posts, or the posts for a given
// tag when the `tag` query param is present.
// @Summary Explore people and posts
// @Tags Explore
// @Produce json
// @Param tag query string false "Filter posts by hashtag"
// @Success 200 {object} exploreResponse
// @Router /explore [get]
func Explore(c *fiber.Ctx) error {
	var (
		ctx      = c.UserContext()
		repo     = repositories.InitializeRepoInstance(ctx)
		viewerID = viewer(c)
	)

	if tag := strings.ToLower(strings.TrimSpace(c.Query("tag"))); tag != "" {
		result, err := repo.FindAll(ctx, &entities.Post{}, &entities.QueryParams{
			Query: entities.Query{
				Joins:   "JOIN post_hashtags ph ON ph.post_id = posts.id",
				Filters: "ph.tag = ?",
				Values:  []any{tag},
			},
			SelectFields: []string{"posts.*"},
			Sort:         "posts.created_at desc",
		})
		if err != nil {
			return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
		}
		return utils.BuildSuccessResponse(c, http.StatusOK, models.DecoratePosts(result, viewerID))
	}

	// Suggested users — newest first, excluding the viewer.
	usersResult, err := repo.FindAll(ctx, &entities.User{}, &entities.QueryParams{
		Query: entities.Query{Filters: "id <> ?", Values: []any{viewerID}},
		Sort:  "created_at desc",
		Limit: suggestedUsersLimit,
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}
	users, _ := usersResult.([]*entities.User)
	if users == nil {
		users = []*entities.User{}
	}

	// Popular posts — decorate all, then keep the most-liked.
	postsResult, err := repo.FindAll(ctx, &entities.Post{}, &entities.QueryParams{Sort: "created_at desc"})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}
	posts := models.TopByLikes(models.DecoratePosts(postsResult, viewerID), popularPostsLimit)

	return utils.BuildSuccessResponse(c, http.StatusOK, exploreResponse{
		Users: users,
		Posts: posts,
	})
}
