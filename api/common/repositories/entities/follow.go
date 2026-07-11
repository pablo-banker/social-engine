package entities

import "github.com/google/uuid"

var _ IEntity = (*Follow)(nil)

// Follow records that one user follows another. It maps to the "follows" table
// and is unique per (follower_id, following_id).
type Follow struct {
	Base
	FollowerID  uuid.UUID `json:"followerId" gorm:"column:follower_id"`
	FollowingID uuid.UUID `json:"followingId" gorm:"column:following_id"`
}

func (f *Follow) TableName() string { return "follows" }

func (f *Follow) LoadAssociations() []string { return nil }

func (f *Follow) GetID() any { return f.ID }
