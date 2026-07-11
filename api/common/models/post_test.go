package models

import (
	"testing"
	"time"

	"social-engine/common/repositories/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ExtractHashtags(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    []string
	}{
		{"basic", "hello #go #svelte", []string{"go", "svelte"}},
		{"dedup and lowercase", "#Go #GO #go", []string{"go"}},
		{"preserves first-seen order", "#b #a #b #c", []string{"b", "a", "c"}},
		{"digits and underscores", "#go_1 #v2", []string{"go_1", "v2"}},
		{"unicode letters", "#café #日本", []string{"café", "日本"}},
		{"ignores bare @ and punctuation", "email@x.com, no tags", []string{}},
		{"trailing punctuation excluded", "read #go! now", []string{"go"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ExtractHashtags(tc.content)
			assert.Equal(t, tc.want, got)
			assert.NotNil(t, got, "result must never be nil")
		})
	}
}

func Test_DecoratePosts(t *testing.T) {
	viewerID := uuid.New()

	t.Run("decorates each post for the viewer", func(t *testing.T) {
		posts := []*entities.Post{
			{
				Base:     entities.Base{ID: uuid.New()},
				Likes:    []entities.Like{{UserID: viewerID}, {UserID: uuid.New()}},
				Comments: []entities.Comment{{}, {}, {}},
			},
			{
				Base:  entities.Base{ID: uuid.New()},
				Likes: []entities.Like{{UserID: uuid.New()}},
			},
		}

		got := DecoratePosts([]*entities.Post{posts[0], posts[1]}, viewerID)
		assert.Len(t, got, 2)

		assert.Equal(t, 2, got[0].LikesCount)
		assert.Equal(t, 3, got[0].CommentsCount)
		assert.True(t, got[0].LikedByMe, "viewer liked the first post")

		assert.Equal(t, 1, got[1].LikesCount)
		assert.Equal(t, 0, got[1].CommentsCount)
		assert.False(t, got[1].LikedByMe, "viewer did not like the second post")
	})

	t.Run("nil result yields empty non-nil slice", func(t *testing.T) {
		got := DecoratePosts(nil, viewerID)
		assert.NotNil(t, got)
		assert.Len(t, got, 0)
	})

	t.Run("wrong type yields empty non-nil slice", func(t *testing.T) {
		got := DecoratePosts("not a post slice", viewerID)
		assert.NotNil(t, got)
		assert.Len(t, got, 0)
	})
}

func Test_TopByLikes(t *testing.T) {
	older := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	newPost := func(likes int, created time.Time) *entities.Post {
		return &entities.Post{
			Base:       entities.Base{ID: uuid.New(), CreatedAt: created},
			LikesCount: likes,
		}
	}

	t.Run("sorts by like count desc and caps at n", func(t *testing.T) {
		a := newPost(1, older)
		b := newPost(3, older)
		c := newPost(2, older)
		input := []*entities.Post{a, b, c}

		got := TopByLikes(input, 2)
		assert.Len(t, got, 2)
		assert.Equal(t, b, got[0])
		assert.Equal(t, c, got[1])

		// input slice order is preserved (not mutated).
		assert.Equal(t, []*entities.Post{a, b, c}, input)
	})

	t.Run("breaks ties by recency", func(t *testing.T) {
		old := newPost(5, older)
		recent := newPost(5, newer)

		got := TopByLikes([]*entities.Post{old, recent}, 0)
		assert.Equal(t, recent, got[0], "more recent post wins the tie")
		assert.Equal(t, old, got[1])
	})

	t.Run("n<=0 returns all sorted", func(t *testing.T) {
		got := TopByLikes([]*entities.Post{newPost(1, older), newPost(9, older)}, 0)
		assert.Len(t, got, 2)
		assert.Equal(t, 9, got[0].LikesCount)
	})

	t.Run("empty input", func(t *testing.T) {
		got := TopByLikes([]*entities.Post{}, 5)
		assert.Len(t, got, 0)
	})
}
