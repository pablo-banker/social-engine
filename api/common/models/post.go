package models

import (
	"regexp"
	"sort"
	"strings"

	"social-engine/common/repositories/entities"

	"github.com/google/uuid"
)

var hashtagRe = regexp.MustCompile(`#([\p{L}\p{N}_]+)`)

// ExtractHashtags returns the unique, lowercased #tags found in the content
// (without the leading '#'), preserving first-seen order.
func ExtractHashtags(content string) []string {
	matches := hashtagRe.FindAllStringSubmatch(content, -1)

	seen := make(map[string]struct{}, len(matches))
	tags := make([]string, 0, len(matches))
	for _, m := range matches {
		tag := strings.ToLower(m[1])
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		tags = append(tags, tag)
	}
	return tags
}

// DecoratePosts type-asserts a repository FindAll result into a posts slice and
// decorates each post for the given viewer (like/comment counts, likedByMe),
// always returning a non-nil slice so JSON responses never serialize as null.
func DecoratePosts(result any, viewerID uuid.UUID) []*entities.Post {
	posts, _ := result.([]*entities.Post)

	decorated := make([]*entities.Post, 0, len(posts))
	for _, post := range posts {
		post.Decorate(viewerID)
		decorated = append(decorated, post)
	}
	return decorated
}

// TopByLikes returns the posts sorted by like count (then recency), capped at n.
// Posts must already be decorated. The input slice is not mutated.
func TopByLikes(posts []*entities.Post, n int) []*entities.Post {
	sorted := make([]*entities.Post, len(posts))
	copy(sorted, posts)

	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].LikesCount != sorted[j].LikesCount {
			return sorted[i].LikesCount > sorted[j].LikesCount
		}
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})

	if n > 0 && len(sorted) > n {
		sorted = sorted[:n]
	}
	return sorted
}
