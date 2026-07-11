package handlers

import (
	"net/http"
	"testing"

	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Explore_Default(t *testing.T) {
	user := &entities.User{Base: entities.Base{ID: uuid.New()}, Name: "Ada", Username: "ada", Bio: "math", AvatarID: "a1"}
	author := entities.User{Base: entities.Base{ID: uuid.New()}, Name: "Ada", Username: "ada", AvatarID: "a1"}
	post := &entities.Post{
		Base:     entities.Base{ID: uuid.New(), CreatedAt: fixedTime},
		AuthorID: author.ID,
		Content:  "popular",
		Author:   author,
		Likes:    []entities.Like{{}, {}},
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler: Explore,
		MockData: []repositories.MockPayload{
			{
				Type: constants.RepositoryFindAll,
				Params: &entities.QueryParams{
					Query: entities.Query{Filters: "id <> ?", Values: []any{uuid.Nil}},
					Sort:  "created_at desc",
					Limit: suggestedUsersLimit,
				},
				ExpectedResult: []*entities.User{user},
			},
			{
				Type:           constants.RepositoryFindAll,
				Params:         &entities.QueryParams{Sort: "created_at desc"},
				ExpectedResult: []*entities.Post{post},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	data := body["data"].(map[string]any)
	assert.Len(t, data["users"].([]any), 1)
	assert.Len(t, data["posts"].([]any), 1)
	assert.Equal(t, "ada", data["users"].([]any)[0].(map[string]any)["username"])
}

func Test_Explore_ByTag(t *testing.T) {
	author := entities.User{Base: entities.Base{ID: uuid.New()}, Name: "Ada", Username: "ada", AvatarID: "a1"}
	post := &entities.Post{
		Base:     entities.Base{ID: uuid.New(), CreatedAt: fixedTime},
		AuthorID: author.ID,
		Content:  "tagged #go",
		Author:   author,
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler:     Explore,
		QueryParams: map[string]any{"tag": "go"},
		MockData: []repositories.MockPayload{
			{
				Type: constants.RepositoryFindAll,
				Params: &entities.QueryParams{
					Query: entities.Query{
						Joins:   "JOIN post_hashtags ph ON ph.post_id = posts.id",
						Filters: "ph.tag = ?",
						Values:  []any{"go"},
					},
					SelectFields: []string{"posts.*"},
					Sort:         "posts.created_at desc",
				},
				ExpectedResult: []*entities.Post{post},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	data, ok := body["data"].([]any)
	assert.True(t, ok, "tag results should be a posts array")
	assert.Len(t, data, 1)
	assert.Equal(t, "tagged #go", data[0].(map[string]any)["content"])
}
