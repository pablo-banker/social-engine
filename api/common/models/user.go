package models

import (
	"context"
	"fmt"
	"strings"

	"social-engine/common/repositories/entities"
	"social-engine/common/repositories/interfaces"

	"github.com/google/uuid"
)

// Defaults for a freshly registered account (mirror the web appearance catalog).
const (
	DefaultAvatarID = "a1"
	DefaultBannerID = "b1"
)

// Appearance catalog — must stay in sync with the web's `lib/appearance.ts`
// (10 avatars a1..a10, 10 banners b1..b10).
var (
	validAvatarIDs = idSet("a", 10)
	validBannerIDs = idSet("b", 10)
)

func idSet(prefix string, n int) map[string]struct{} {
	set := make(map[string]struct{}, n)
	for i := 1; i <= n; i++ {
		set[fmt.Sprintf("%s%d", prefix, i)] = struct{}{}
	}
	return set
}

// IsValidAvatarID reports whether id is one of the catalog avatars.
func IsValidAvatarID(id string) bool {
	_, ok := validAvatarIDs[id]
	return ok
}

// IsValidBannerID reports whether id is one of the catalog banners.
func IsValidBannerID(id string) bool {
	_, ok := validBannerIDs[id]
	return ok
}

// UsernameFromEmail turns an email into a safe, lowercased username handle.
func UsernameFromEmail(email string) string {
	local, _, _ := strings.Cut(email, "@")

	var b strings.Builder
	for _, r := range strings.ToLower(local) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_', r == '.':
			b.WriteRune(r)
		}
	}

	username := strings.Trim(b.String(), ".")
	if username == "" {
		username = "user"
	}
	return username
}

// UniqueUsername derives a username from the email local-part and appends a
// short suffix on collision. The unique index on users.username is the final
// guard, so a single existence check is enough in practice.
func UniqueUsername(ctx context.Context, repo interfaces.IRepository, email string) (string, error) {
	base := UsernameFromEmail(email)

	taken, err := repo.Verify(ctx, &entities.User{}, &entities.QueryParams{
		Query: entities.Query{Filters: "lower(username) = ?", Values: []any{base}},
	})
	if err != nil {
		return "", err
	}
	if !taken {
		return base, nil
	}

	return base + "-" + uuid.NewString()[:6], nil
}
