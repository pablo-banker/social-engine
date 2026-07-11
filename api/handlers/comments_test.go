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

func Test_ListComments(t *testing.T) {
	postID := uuid.New()
	author := entities.User{Base: entities.Base{ID: uuid.New()}, Name: "Ada", Username: "ada", AvatarID: "a1"}
	comment := &entities.Comment{
		Base:     entities.Base{ID: uuid.New(), CreatedAt: fixedTime},
		PostID:   postID,
		AuthorID: author.ID,
		Content:  "nice post",
		Author:   author,
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler:   ListComments,
		URLParams: map[string]string{"id": postID.String()},
		MockData: []repositories.MockPayload{
			{
				Type: constants.RepositoryFindAll,
				Params: &entities.QueryParams{
					Query: entities.Query{Filters: "post_id = ?", Values: []any{postID}},
					Sort:  "created_at asc",
				},
				ExpectedResult: []*entities.Comment{comment},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	data := body["data"].([]any)
	assert.Len(t, data, 1)
	assert.Equal(t, "nice post", data[0].(map[string]any)["content"])
	assert.Equal(t, "ada", data[0].(map[string]any)["author"].(map[string]any)["username"])
}

func Test_AddComment(t *testing.T) {
	authorID := uuid.New()
	postID := uuid.New()

	postFilter := &entities.QueryParams{
		Query: entities.Query{Filters: "id = ?", Values: []any{postID}},
	}

	t.Run("success", func(t *testing.T) {
		commentID := uuid.New()
		saved := entities.Comment{
			Base:     entities.Base{ID: commentID, CreatedAt: fixedTime},
			PostID:   postID,
			AuthorID: authorID,
			Content:  "great write-up",
		}
		author := entities.User{Base: entities.Base{ID: authorID}, Name: "Ada", Username: "ada", AvatarID: "a1"}
		reloaded := entities.Comment{
			Base:     entities.Base{ID: commentID, CreatedAt: fixedTime},
			PostID:   postID,
			AuthorID: authorID,
			Content:  "great write-up",
			Author:   author,
		}

		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:     AddComment,
			URLParams:   map[string]string{"id": postID.String()},
			Body:        createCommentRequest{Content: "great write-up"},
			Middlewares: []fiber.Handler{withUser(authorID)},
			MockData: []repositories.MockPayload{
				{Type: constants.RepositoryVerify, Params: postFilter, ExpectedResult: true},
				{Type: constants.RepositorySave, ExpectedResult: saved},
				{Type: constants.RepositoryFindByID, ExpectedResult: reloaded},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, statusCode)

		data := body["data"].(map[string]any)
		assert.Equal(t, "great write-up", data["content"])
		assert.Equal(t, "ada", data["author"].(map[string]any)["username"])
	})

	t.Run("post not found", func(t *testing.T) {
		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:     AddComment,
			URLParams:   map[string]string{"id": postID.String()},
			Body:        createCommentRequest{Content: "great write-up"},
			Middlewares: []fiber.Handler{withUser(authorID)},
			MockData: []repositories.MockPayload{
				{Type: constants.RepositoryVerify, Params: postFilter, ExpectedResult: false},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, apiErrors.ErrPostNotFound.Code, body["code"])
	})
}
