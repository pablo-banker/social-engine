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
	"gorm.io/gorm"
)

func Test_ListFeed(t *testing.T) {
	viewerID := uuid.New()
	author := entities.User{Base: entities.Base{ID: uuid.New()}, Name: "Ada", Username: "ada", AvatarID: "a1"}
	post := &entities.Post{
		Base:     entities.Base{ID: uuid.New(), CreatedAt: fixedTime},
		AuthorID: author.ID,
		Content:  "hello #go",
		Author:   author,
		Likes:    []entities.Like{{UserID: viewerID}, {UserID: uuid.New()}},
		Comments: []entities.Comment{{Content: "nice"}},
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler:     List,
		Middlewares: []fiber.Handler{withUser(viewerID)},
		MockData: []repositories.MockPayload{
			{
				Type:           constants.RepositoryFindAll,
				Params:         &entities.QueryParams{Sort: "created_at desc"},
				ExpectedResult: []*entities.Post{post},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	data, ok := body["data"].([]any)
	assert.True(t, ok, "data should be an array")
	assert.Len(t, data, 1)

	p := data[0].(map[string]any)
	assert.Equal(t, "hello #go", p["content"])
	assert.Equal(t, float64(2), p["likes"])
	assert.Equal(t, float64(1), p["comments"])
	assert.Equal(t, true, p["likedByMe"])
	assert.Equal(t, "ada", p["author"].(map[string]any)["username"])
}

func Test_GetPost(t *testing.T) {
	postID := uuid.New()

	t.Run("found", func(t *testing.T) {
		author := entities.User{Base: entities.Base{ID: uuid.New()}, Name: "Ada", Username: "ada", AvatarID: "a1"}
		post := entities.Post{
			Base:     entities.Base{ID: postID, CreatedAt: fixedTime},
			AuthorID: author.ID,
			Content:  "single post",
			Author:   author,
		}

		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:   Get,
			URLParams: map[string]string{"id": postID.String()},
			MockData:  []repositories.MockPayload{{Type: constants.RepositoryFindByID, ExpectedResult: post}},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, statusCode)

		data := body["data"].(map[string]any)
		assert.Equal(t, "single post", data["content"])
		assert.Equal(t, "ada", data["author"].(map[string]any)["username"])
	})

	t.Run("not found", func(t *testing.T) {
		statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
			Handler:   Get,
			URLParams: map[string]string{"id": postID.String()},
			MockData:  []repositories.MockPayload{{Type: constants.RepositoryFindByID, ExpectedError: gorm.ErrRecordNotFound}},
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, statusCode)
		assert.Equal(t, apiErrors.ErrPostNotFound.Code, body["code"])
	})
}

func Test_CreatePost(t *testing.T) {
	authorID := uuid.New()
	postID := uuid.New()

	savedPost := entities.Post{
		Base:     entities.Base{ID: postID, CreatedAt: fixedTime},
		AuthorID: authorID,
		Content:  "hello #go #svelte",
	}
	author := entities.User{Base: entities.Base{ID: authorID}, Name: "Ada Lovelace", Username: "ada", AvatarID: "a1"}
	reloaded := entities.Post{
		Base:     entities.Base{ID: postID, CreatedAt: fixedTime},
		AuthorID: authorID,
		Content:  "hello #go #svelte",
		Author:   author,
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler:     Create,
		Body:        createPostRequest{Content: "hello #go #svelte"},
		Middlewares: []fiber.Handler{withUser(authorID)},
		MockData: []repositories.MockPayload{
			{Type: constants.RepositoryBeginTx},
			{Type: constants.RepositorySave, ExpectedResult: savedPost},
			{Type: constants.RepositoryBulkSave},
			{Type: constants.RepositoryCommit},
			{Type: constants.RepositoryFindByID, ExpectedResult: reloaded},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, statusCode)

	data := body["data"].(map[string]any)
	assert.Equal(t, "hello #go #svelte", data["content"])
	assert.Equal(t, float64(0), data["likes"])
	assert.Equal(t, "ada", data["author"].(map[string]any)["username"])
}

func Test_CreatePost_ValidationError(t *testing.T) {
	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler:     Create,
		Body:        createPostRequest{Content: "   "},
		Middlewares: []fiber.Handler{withUser(uuid.New())},
		MockData:    []repositories.MockPayload{},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, apiErrors.ErrValidation.Code, body["code"])
}
