package models

import (
	"context"

	"social-engine/common/repositories/customs"
	"social-engine/common/repositories/interfaces"
)

// TrendingTopic is a hashtag with its post count. It is a GROUP BY aggregate
// over post_hashtags, so it has no backing entity and stays a response DTO.
type TrendingTopic struct {
	Tag   string `json:"tag"`
	Posts int    `json:"posts"`
}

// TrendingTopics returns the hottest hashtags with their post counts, capped at
// limit. The result is always a non-nil slice.
func TrendingTopics(ctx context.Context, repo interfaces.IRepository, limit int) ([]TrendingTopic, error) {
	topics := []TrendingTopic{}
	if err := repo.Raw(ctx, &topics, customs.TrendingTopicsQuery(), limit); err != nil {
		return nil, err
	}
	return topics, nil
}
