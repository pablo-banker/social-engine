package handlers

import (
	"errors"
	"net/http"
	"strings"

	"social-engine/common/apiErrors"
	"social-engine/common/models"
	"social-engine/common/repositories"
	"social-engine/common/repositories/entities"
	"social-engine/common/repositories/interfaces"
	"social-engine/common/utils"
	"social-engine/common/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type createPostRequest struct {
	Content string `json:"content" validate:"required,max=500"`
}

// List returns the public feed, newest first.
// @Summary List the public feed
// @Description Public feed of posts; likedByMe reflects the caller when authenticated
// @Tags Posts
// @Produce json
// @Success 200 {array} entities.Post
// @Router /posts [get]
func List(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	result, err := repo.FindAll(ctx, &entities.Post{}, &entities.QueryParams{Sort: "created_at desc"})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, models.DecoratePosts(result, viewer(c)))
}

// Get returns a single decorated post.
// @Summary Get a post
// @Tags Posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} entities.Post
// @Failure 404 {object} utils.ErrorResponse
// @Router /posts/{id} [get]
func Get(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	postID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrPostNotFound)
	}

	var post entities.Post
	if err := repo.FindByID(ctx, &post, postID, nil); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.BuildErrorResponse(c, apiErrors.ErrPostNotFound)
		}
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	post.Decorate(viewer(c))
	return utils.BuildSuccessResponse(c, http.StatusOK, &post)
}

// Create publishes a new post and extracts its #hashtags.
// @Summary Create a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param payload body createPostRequest true "Post content"
// @Success 201 {object} entities.Post
// @Failure 401 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /posts [post]
func Create(c *fiber.Ctx) error {
	var (
		ctx      = c.UserContext()
		repo     = repositories.InitializeRepoInstance(ctx)
		authorID = viewer(c)
	)

	var req createPostRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidPayload.WithInternal(err.Error()))
	}

	req.Content = strings.TrimSpace(req.Content)

	if err := validation.Struct(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrValidation.WithDetails(err.Error()))
	}

	post := entities.Post{AuthorID: authorID, Content: req.Content}
	tags := models.ExtractHashtags(req.Content)

	err := repo.WithTransaction(ctx, func(tx interfaces.IRepository) error {
		if err := tx.Save(ctx, &post, nil); err != nil {
			return err
		}

		if len(tags) > 0 {
			hashtags := make([]entities.IEntity, 0, len(tags))
			for _, tag := range tags {
				hashtags = append(hashtags, &entities.PostHashtag{PostID: post.ID, Tag: tag})
			}
			if err := tx.BulkSave(ctx, hashtags); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	// Reload with associations so the response carries the nested author.
	if err := repo.FindByID(ctx, &post, post.ID, nil); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	post.Decorate(authorID)
	return utils.BuildSuccessResponse(c, http.StatusCreated, &post)
}
