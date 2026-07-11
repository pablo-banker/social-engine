package entities

import "github.com/google/uuid"

var _ IEntity = (*PostHashtag)(nil)

// PostHashtag is a single #tag extracted from a post's content on creation.
// It maps to the "post_hashtags" table and powers trending and explore-by-tag.
type PostHashtag struct {
	Base
	PostID uuid.UUID `json:"postId" gorm:"column:post_id"`
	Tag    string    `json:"tag" gorm:"column:tag"`
}

func (h *PostHashtag) TableName() string { return "post_hashtags" }

func (h *PostHashtag) LoadAssociations() []string { return nil }

func (h *PostHashtag) GetID() any { return h.ID }
