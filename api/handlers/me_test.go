package handlers

import (
	"net/http"
	"testing"

	"social-engine/common/apiErrors"
	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetSettings(t *testing.T) {
	userID := uuid.New()
	user := entities.User{
		Base: entities.Base{ID: userID}, Name: "Ada", Bio: "math", AvatarID: "a1", BannerID: "b1",
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler:     GetSettings,
		Middlewares: []fiber.Handler{withUser(userID)},
		MockData:    []repositories.MockPayload{{Type: constants.RepositoryFindByID, ExpectedResult: user}},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	data := body["data"].(map[string]any)
	assert.Equal(t, "Ada", data["name"])
	assert.Equal(t, "a1", data["avatarId"])
	assert.Equal(t, "b1", data["bannerId"])
}

func Test_UpdateProfile(t *testing.T) {
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		updated := &entities.User{
			Base:     entities.Base{ID: userID},
			Name:     "Ada L.",
			Username: "ada",
			Email:    "ada@example.com",
			Bio:      "countess of lovelace",
			AvatarID: "a3",
			BannerID: "b5",
		}

		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:     UpdateProfile,
			Middlewares: []fiber.Handler{withUser(userID)},
			Body: updateProfileRequest{
				Name: "Ada L.", Bio: "countess of lovelace", AvatarID: "a3", BannerID: "b5",
			},
			MockData: []repositories.MockPayload{
				{
					Type:           constants.RepositoryUpdate,
					Params:         &entities.QueryParams{UpdateFields: []string{"name", "bio", "avatar_id", "banner_id"}},
					ExpectedResult: updated,
				},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)

		data := body["data"].(map[string]any)
		assert.Equal(t, "ada", data["username"])
		assert.Equal(t, "ada@example.com", data["email"])
		assert.Equal(t, "a3", data["avatarId"])
	})

	t.Run("invalid avatar", func(t *testing.T) {
		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:     UpdateProfile,
			Middlewares: []fiber.Handler{withUser(userID)},
			Body: updateProfileRequest{
				Name: "Ada", Bio: "", AvatarID: "z99", BannerID: "b1",
			},
			MockData: []repositories.MockPayload{},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
		assert.Equal(t, apiErrors.ErrInvalidAvatar.Code, body["code"])
	})
}
