package entities

import "github.com/google/uuid"

var _ IEntity = (*Post)(nil)

// Post is a message published by a user. It maps to the "posts" table.
type Post struct {
	Base
	AuthorID uuid.UUID `json:"authorId" gorm:"column:author_id"`
	Content  string    `json:"content" gorm:"column:content"`

	// Associations loaded via LoadAssociations. Author is exposed nested;
	// Likes/Comments are loaded only to derive the counters below.
	Author   User      `json:"author" gorm:"foreignKey:AuthorID;references:ID"`
	Likes    []Like    `json:"-" gorm:"foreignKey:PostID;references:ID"`
	Comments []Comment `json:"-" gorm:"foreignKey:PostID;references:ID"`

	// Transient view fields derived from the loaded associations.
	LikesCount    int  `json:"likes" gorm:"-"`
	CommentsCount int  `json:"comments" gorm:"-"`
	LikedByMe     bool `json:"likedByMe" gorm:"-"`
}

func (p *Post) TableName() string { return "posts" }

func (p *Post) LoadAssociations() []string {
	return []string{"Author", "Likes", "Comments"}
}

func (p *Post) GetID() any { return p.ID }

// Decorate fills the transient view fields from the loaded associations,
// flagging whether the given viewer liked the post.
func (p *Post) Decorate(viewerID uuid.UUID) {
	p.LikesCount = len(p.Likes)
	p.CommentsCount = len(p.Comments)
	for _, like := range p.Likes {
		if like.UserID == viewerID {
			p.LikedByMe = true
			break
		}
	}
}
