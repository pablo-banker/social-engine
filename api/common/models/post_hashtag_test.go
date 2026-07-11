package models

import (
	"context"
	"errors"
	"testing"

	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"

	"github.com/stretchr/testify/assert"
)

func Test_TrendingTopics(t *testing.T) {
	t.Run("returns the aggregated topics", func(t *testing.T) {
		want := []TrendingTopic{{Tag: "go", Posts: 12}, {Tag: "svelte", Posts: 7}}
		repo := repositories.NewMockRepository([]repositories.MockPayload{
			{Type: constants.RepositoryRaw, ExpectedResult: want},
		})

		got, err := TrendingTopics(context.Background(), repo, 8)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		repo := repositories.NewMockRepository([]repositories.MockPayload{
			{Type: constants.RepositoryRaw, ExpectedError: errors.New("boom")},
		})

		got, err := TrendingTopics(context.Background(), repo, 8)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
