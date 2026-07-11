package handlers

import (
	"net/http"

	"social-engine/common/apiErrors"
	"social-engine/common/repositories"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// likeResult is the outcome of toggling a like. It is a purpose-built action
// response with no backing entity.
type likeResult struct {
	Liked bool `json:"liked"`
	Likes int  `json:"likes"`
}

// ToggleLike likes the post if not already liked, otherwise unlikes it, and
// returns the new like state and total.
// @Summary Toggle a like on a post
// @Tags Likes
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} likeResult
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/like [post]
func ToggleLike(c *fiber.Ctx) error {
	var (
		ctx    = c.UserContext()
		repo   = repositories.InitializeRepoInstance(ctx)
		userID = viewer(c)
	)

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrPostNotFound)
	}

	postExists, err := repo.Verify(ctx, &entities.Post{}, &entities.QueryParams{
		Query: entities.Query{Filters: "id = ?", Values: []any{postID}},
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}
	if !postExists {
		return utils.BuildErrorResponse(c, apiErrors.ErrPostNotFound)
	}

	likeFilter := entities.Query{Filters: "user_id = ? AND post_id = ?", Values: []any{userID, postID}}

	liked, err := repo.Verify(ctx, &entities.Like{}, &entities.QueryParams{Query: likeFilter})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	if liked {
		if err := repo.Delete(ctx, &entities.Like{}, &entities.QueryParams{Query: likeFilter}); err != nil {
			return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
		}
	} else {
		if err := repo.Save(ctx, &entities.Like{UserID: userID, PostID: postID}, nil); err != nil {
			return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
		}
	}

	count, err := repo.Count(ctx, &entities.Like{}, &entities.QueryParams{
		Query: entities.Query{Filters: "post_id = ?", Values: []any{postID}},
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, likeResult{
		Liked: !liked,
		Likes: int(count),
	})
}
