package handlers

import (
	"net/http"
	"testing"

	"social-engine/common/models"
	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Trending(t *testing.T) {
	topics := []models.TrendingTopic{
		{Tag: "go", Posts: 12},
		{Tag: "svelte", Posts: 7},
	}
	author := entities.User{Base: entities.Base{ID: uuid.New()}, Name: "Ada", Username: "ada", AvatarID: "a1"}
	post := &entities.Post{
		Base:     entities.Base{ID: uuid.New(), CreatedAt: fixedTime},
		AuthorID: author.ID,
		Content:  "popular",
		Author:   author,
		Likes:    []entities.Like{{}, {}},
	}

	statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
		Handler: Trending,
		MockData: []repositories.MockPayload{
			{Type: constants.RepositoryRaw, ExpectedResult: topics},
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
	topicsData := data["topics"].([]any)
	assert.Len(t, topicsData, 2)
	assert.Equal(t, "go", topicsData[0].(map[string]any)["tag"])
	assert.Equal(t, float64(12), topicsData[0].(map[string]any)["posts"])
	assert.Len(t, data["posts"].([]any), 1)
}
