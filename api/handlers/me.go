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
	"gorm.io/gorm"
)

type updateProfileRequest struct {
	Name     string `json:"name" validate:"required,max=50"`
	Bio      string `json:"bio" validate:"max=160"`
	AvatarID string `json:"avatarId" validate:"required"`
	BannerID string `json:"bannerId" validate:"required"`
}

// GetSettings returns the logged-in user's editable profile fields.
// @Summary Get my profile settings
// @Tags Me
// @Produce json
// @Success 200 {object} entities.User
// @Failure 401 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /me [get]
func GetSettings(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	var user entities.User
	if err := repo.FindByID(ctx, &user, viewer(c), nil); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.BuildErrorResponse(c, apiErrors.ErrUserNotFound)
		}
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, &user)
}

// UpdateProfile updates the logged-in user's name, bio, avatar and banner.
// @Summary Update my profile
// @Tags Me
// @Accept json
// @Produce json
// @Param payload body updateProfileRequest true "Profile fields"
// @Success 200 {object} userDTO
// @Failure 401 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /me [patch]
func UpdateProfile(c *fiber.Ctx) error {
	var (
		ctx    = c.UserContext()
		repo   = repositories.InitializeRepoInstance(ctx)
		userID = viewer(c)
	)

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidPayload.WithInternal(err.Error()))
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Bio = strings.TrimSpace(req.Bio)

	if err := validation.Struct(&req); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrValidation.WithDetails(err.Error()))
	}
	if !models.IsValidAvatarID(req.AvatarID) {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidAvatar)
	}
	if !models.IsValidBannerID(req.BannerID) {
		return utils.BuildErrorResponse(c, apiErrors.ErrInvalidBanner)
	}

	user := entities.User{
		Base:     entities.Base{ID: userID},
		Name:     req.Name,
		Bio:      req.Bio,
		AvatarID: req.AvatarID,
		BannerID: req.BannerID,
	}
	if err := repo.Update(ctx, &user, &entities.QueryParams{
		UpdateFields: []string{"name", "bio", "avatar_id", "banner_id"},
	}); err != nil {
		return utils.BuildErrorResponse(c, apiErrors.ErrInternal.WithInternal(err.Error()))
	}

	return utils.BuildSuccessResponse(c, http.StatusOK, userDTO{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		AvatarID: user.AvatarID,
	})
}
