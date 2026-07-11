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
	"social-engine/common/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type registerRequest struct {
	FirstName string `json:"firstName" validate:"required,max=50"`
	LastName  string `json:"lastName" validate:"required,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type userDTO struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	AvatarID string    `json:"avatarId"`
}

type authResponse struct {
	User        userDTO `json:"user"`
	AccessToken string  `json:"accessToken"`
}

// Register creates a new account and returns the user with an access token.
// @Summary Register a new account
// @Description Create an account and receive an access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body registerRequest true "Registration data"
// @Success 201 {object} authResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidPayload.WithInternal(err.Error()))
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validation.Struct(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrValidation.WithDetails(err.Error()))
	}

	taken, err := repo.Verify(ctx, &entities.User{}, &entities.QueryParams{
		Query: entities.Query{Filters: "lower(email) = ?", Values: []any{req.Email}},
	})
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}
	if taken {
		return utils.BuildErrorResponse(c, apiErrors.ErrEmailTaken)
	}

	username, err := models.UniqueUsername(ctx, repo, req.Email)
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	hash, err := models.HashPassword(req.Password)
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	user := entities.User{
		Name:         strings.TrimSpace(req.FirstName + " " + req.LastName),
		Username:     username,
		Email:        req.Email,
		PasswordHash: hash,
		AvatarID:     models.DefaultAvatarID,
		BannerID:     models.DefaultBannerID,
	}

	if err := repo.Save(ctx, &user, nil); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	token, err := models.GenerateToken(user.ID)
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusCreated, authResponse{
		User: userDTO{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			AvatarID: user.AvatarID,
		},
		AccessToken: token,
	})
}

// Login authenticates by email/password and returns an access token.
// @Summary Log in
// @Description Authenticate with email and password and receive an access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body loginRequest true "Login credentials"
// @Success 200 {object} authResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidPayload.WithInternal(err.Error()))
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validation.Struct(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrValidation.WithDetails(err.Error()))
	}

	var user entities.User
	err := repo.FindOne(ctx, &user, &entities.QueryParams{
		Query: entities.Query{Filters: "lower(email) = ?", Values: []any{req.Email}},
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.BuildErrorResponse(c, apiErrors.ErrInvalidCredentials)
		}
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	if !models.CheckPassword(user.PasswordHash, req.Password) {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidCredentials)
	}

	token, err := models.GenerateToken(user.ID)
	if err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, authResponse{
		User: userDTO{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			AvatarID: user.AvatarID,
		},
		AccessToken: token,
	})
}
