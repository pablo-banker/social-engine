package handlers

import (
	"net/http"
	"testing"

	"social-engine/common/apiErrors"
	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_GetProfile(t *testing.T) {
	usernameFilter := func(username string) *entities.QueryParams {
		return &entities.QueryParams{
			Query: entities.Query{Filters: "lower(username) = ?", Values: []any{username}},
		}
	}

	t.Run("found", func(t *testing.T) {
		userID := uuid.New()
		user := entities.User{
			Base:     entities.Base{ID: userID, CreatedAt: fixedTime},
			Name:     "Ada Lovelace",
			Username: "ada",
			Bio:      "math",
			AvatarID: "a1",
			BannerID: "b1",
		}

		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:   GetProfile,
			URLParams: map[string]string{"username": "ada"},
			MockData: []repositories.MockPayload{
				{Type: constants.RepositoryFindOne, Params: usernameFilter("ada"), ExpectedResult: user},
				{Type: constants.RepositoryCount, Params: &entities.QueryParams{
					Query: entities.Query{Filters: "author_id = ?", Values: []any{userID}},
				}, ExpectedResult: 5},
				{Type: constants.RepositoryCount, Params: &entities.QueryParams{
					Query: entities.Query{Filters: "following_id = ?", Values: []any{userID}},
				}, ExpectedResult: 10},
				{Type: constants.RepositoryCount, Params: &entities.QueryParams{
					Query: entities.Query{Filters: "follower_id = ?", Values: []any{userID}},
				}, ExpectedResult: 3},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)

		data := body["data"].(map[string]any)
		assert.Equal(t, "ada", data["username"])
		stats := data["stats"].(map[string]any)
		assert.Equal(t, float64(5), stats["posts"])
		assert.Equal(t, float64(10), stats["followers"])
		assert.Equal(t, float64(3), stats["following"])
	})

	t.Run("not found", func(t *testing.T) {
		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:   GetProfile,
			URLParams: map[string]string{"username": "ghost"},
			MockData: []repositories.MockPayload{
				{Type: constants.RepositoryFindOne, Params: usernameFilter("ghost"), ExpectedError: gorm.ErrRecordNotFound},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, apiErrors.ErrUserNotFound.Code, body["code"])
	})
}

func Test_ListUserPosts(t *testing.T) {
	userID := uuid.New()
	user := entities.User{Base: entities.Base{ID: userID}, Username: "ada"}
	author := entities.User{Base: entities.Base{ID: userID}, Name: "Ada", Username: "ada", AvatarID: "a1"}
	post := &entities.Post{
		Base:     entities.Base{ID: uuid.New(), CreatedAt: fixedTime},
		AuthorID: userID,
		Content:  "my post",
		Author:   author,
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler:   ListUserPosts,
		URLParams: map[string]string{"username": "ada"},
		MockData: []repositories.MockPayload{
			{Type: constants.RepositoryFindOne, Params: &entities.QueryParams{
				Query: entities.Query{Filters: "lower(username) = ?", Values: []any{"ada"}},
			}, ExpectedResult: user},
			{Type: constants.RepositoryFindAll, Params: &entities.QueryParams{
				Query: entities.Query{Filters: "author_id = ?", Values: []any{userID}},
				Sort:  "created_at desc",
			}, ExpectedResult: []*entities.Post{post}},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Len(t, body["data"].([]any), 1)
}
