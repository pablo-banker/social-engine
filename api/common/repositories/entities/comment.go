package entities

import "github.com/google/uuid"

var _ IEntity = (*Comment)(nil)

// Comment is a reply to a post. It maps to the "comments" table.
type Comment struct {
	Base
	PostID   uuid.UUID `json:"postId" gorm:"column:post_id"`
	AuthorID uuid.UUID `json:"authorId" gorm:"column:author_id"`
	Content  string    `json:"content" gorm:"column:content"`

	Author User `json:"author" gorm:"foreignKey:AuthorID;references:ID"`
}

func (c *Comment) TableName() string { return "comments" }

func (c *Comment) LoadAssociations() []string { return []string{"Author"} }

func (c *Comment) GetID() any { return c.ID }
