package models

import (
	"context"
	"strings"
	"testing"

	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/entities"

	"github.com/stretchr/testify/assert"
)

func Test_UsernameFromEmail(t *testing.T) {
	cases := []struct {
		name  string
		email string
		want  string
	}{
		{"simple", "ada@example.com", "ada"},
		{"lowercased", "Ada@Example.com", "ada"},
		{"keeps dots and underscores", "Ada.Lovelace_1@x.com", "ada.lovelace_1"},
		{"strips disallowed chars", "user+spam!@x.com", "userspam"},
		{"digits allowed", "123@x.com", "123"},
		{"trims leading/trailing dots", ".ada.@x.com", "ada"},
		{"empty local falls back", "@x.com", "user"},
		{"dots-only falls back", "...@x.com", "user"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, UsernameFromEmail(tc.email))
		})
	}
}

func Test_UniqueUsername_Available(t *testing.T) {
	repo := repositories.NewMockRepository([]repositories.MockPayload{
		{
			Type: constants.RepositoryVerify,
			Params: &entities.QueryParams{
				Query: entities.Query{Filters: "lower(username) = ?", Values: []any{"ada"}},
			},
			ExpectedResult: false,
		},
	})

	got, err := UniqueUsername(context.Background(), repo, "ada@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "ada", got)
}

func Test_UniqueUsername_Taken(t *testing.T) {
	repo := repositories.NewMockRepository([]repositories.MockPayload{
		{
			Type: constants.RepositoryVerify,
			Params: &entities.QueryParams{
				Query: entities.Query{Filters: "lower(username) = ?", Values: []any{"ada"}},
			},
			ExpectedResult: true,
		},
	})

	got, err := UniqueUsername(context.Background(), repo, "ada@example.com")
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(got, "ada-"), "collision should append a suffix")
	assert.Len(t, got, len("ada-")+6, "suffix should be a 6-char shortid")
}

func Test_IsValidAvatarID(t *testing.T) {
	assert.True(t, IsValidAvatarID("a1"))
	assert.True(t, IsValidAvatarID("a10"))
	assert.True(t, IsValidAvatarID(DefaultAvatarID))
	assert.False(t, IsValidAvatarID("a0"))
	assert.False(t, IsValidAvatarID("a11"))
	assert.False(t, IsValidAvatarID("b1"))
	assert.False(t, IsValidAvatarID(""))
}

func Test_IsValidBannerID(t *testing.T) {
	assert.True(t, IsValidBannerID("b1"))
	assert.True(t, IsValidBannerID("b10"))
	assert.True(t, IsValidBannerID(DefaultBannerID))
	assert.False(t, IsValidBannerID("b0"))
	assert.False(t, IsValidBannerID("b11"))
	assert.False(t, IsValidBannerID("a1"))
	assert.False(t, IsValidBannerID(""))
}
