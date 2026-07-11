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

func Test_ToggleLike(t *testing.T) {
	userID := uuid.New()
	postID := uuid.New()

	postFilter := &entities.QueryParams{
		Query: entities.Query{Filters: "id = ?", Values: []any{postID}},
	}
	likeFilter := &entities.QueryParams{
		Query: entities.Query{Filters: "user_id = ? AND post_id = ?", Values: []any{userID, postID}},
	}
	countFilter := &entities.QueryParams{
		Query: entities.Query{Filters: "post_id = ?", Values: []any{postID}},
	}

	t.Run("like (was not liked)", func(t *testing.T) {
		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:     ToggleLike,
			URLParams:   map[string]string{"id": postID.String()},
			Middlewares: []fiber.Handler{withUser(userID)},
			MockData: []repositories.MockPayload{
				{Type: constants.RepositoryVerify, Params: postFilter, ExpectedResult: true},
				{Type: constants.RepositoryVerify, Params: likeFilter, ExpectedResult: false},
				{Type: constants.RepositorySave, ExpectedResult: entities.Like{Base: entities.Base{ID: uuid.New()}}},
				{Type: constants.RepositoryCount, Params: countFilter, ExpectedResult: 3},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)

		data := body["data"].(map[string]any)
		assert.Equal(t, true, data["liked"])
		assert.Equal(t, float64(3), data["likes"])
	})

	t.Run("unlike (was liked)", func(t *testing.T) {
		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:     ToggleLike,
			URLParams:   map[string]string{"id": postID.String()},
			Middlewares: []fiber.Handler{withUser(userID)},
			MockData: []repositories.MockPayload{
				{Type: constants.RepositoryVerify, Params: postFilter, ExpectedResult: true},
				{Type: constants.RepositoryVerify, Params: likeFilter, ExpectedResult: true},
				{Type: constants.RepositoryDelete, Params: likeFilter},
				{Type: constants.RepositoryCount, Params: countFilter, ExpectedResult: 2},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)

		data := body["data"].(map[string]any)
		assert.Equal(t, false, data["liked"])
		assert.Equal(t, float64(2), data["likes"])
	})

	t.Run("post not found", func(t *testing.T) {
		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:     ToggleLike,
			URLParams:   map[string]string{"id": postID.String()},
			Middlewares: []fiber.Handler{withUser(userID)},
			MockData: []repositories.MockPayload{
				{Type: constants.RepositoryVerify, Params: postFilter, ExpectedResult: false},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, apiErrors.ErrPostNotFound.Code, body["code"])
	})
}
