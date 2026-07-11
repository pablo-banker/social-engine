package entities

import "github.com/google/uuid"

var _ IEntity = (*Like)(nil)

// Like records that a user liked a post. It maps to the "likes" table and is
// unique per (user_id, post_id).
type Like struct {
	Base
	UserID uuid.UUID `json:"userId" gorm:"column:user_id"`
	PostID uuid.UUID `json:"postId" gorm:"column:post_id"`
}

func (l *Like) TableName() string { return "likes" }

func (l *Like) LoadAssociations() []string { return nil }

func (l *Like) GetID() any { return l.ID }
