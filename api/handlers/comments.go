package handlers

import (
	"net/http"
	"strings"

	"social-engine/common/apiErrors"
	"social-engine/common/repositories"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"
	"social-engine/common/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createCommentRequest struct {
	Content string `json:"content" validate:"required,max=300"`
}

// ListComments returns a post's comments, oldest first.
// @Summary List a post's comments
// @Tags Comments
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {array} entities.Comment
// @Router /posts/{id}/comments [get]
func ListComments(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrPostNotFound)
	}

	result, err := repo.FindAll(ctx, &entities.Comment{}, &entities.QueryParams{
		Query: entities.Query{Filters: "post_id = ?", Values: []any{postID}},
		Sort:  "created_at asc",
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	comments, _ := result.([]*entities.Comment)
	if comments == nil {
		comments = []*entities.Comment{}
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, comments)
}

// AddComment publishes a comment on a post.
// @Summary Add a comment
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param payload body createCommentRequest true "Comment content"
// @Success 201 {object} entities.Comment
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /posts/{id}/comments [post]
func AddComment(c *fiber.Ctx) error {
	var (
		ctx      = c.UserContext()
		repo     = repositories.InitializeRepoInstance(ctx)
		authorID = viewer(c)
	)

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrPostNotFound)
	}

	var req createCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidPayload.WithInternal(err.Error()))
	}

	req.Content = strings.TrimSpace(req.Content)

	if err := validation.Struct(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrValidation.WithDetails(err.Error()))
	}

	exists, err := repo.Verify(ctx, &entities.Post{}, &entities.QueryParams{
		Query: entities.Query{Filters: "id = ?", Values: []any{postID}},
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}
	if !exists {
		return utils.BuildErrorResponse(c, apiErrors.ErrPostNotFound)
	}

	comment := entities.Comment{PostID: postID, AuthorID: authorID, Content: req.Content}
	if err := repo.Save(ctx, &comment, nil); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	// Reload with the author association for the nested response.
	if err := repo.FindByID(ctx, &comment, comment.ID, nil); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusCreated, &comment)
}
