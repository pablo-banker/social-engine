package entities

var _ IEntity = (*User)(nil)

// User is a registered account. It maps to the "users" table.
type User struct {
	Base
	Name         string `json:"name" gorm:"column:name"`
	Username     string `json:"username" gorm:"column:username"`
	Email        string `json:"-" gorm:"column:email"`
	PasswordHash string `json:"-" gorm:"column:password_hash"`
	Bio          string `json:"bio" gorm:"column:bio"`
	AvatarID     string `json:"avatarId" gorm:"column:avatar_id"`
	BannerID     string `json:"bannerId" gorm:"column:banner_id"`

	// Stats is a transient aggregate populated only for profile responses.
	Stats *UserStats `json:"stats,omitempty" gorm:"-"`
}

// UserStats holds a profile's aggregate counters.
type UserStats struct {
	Posts     int `json:"posts"`
	Followers int `json:"followers"`
	Following int `json:"following"`
}

func (u *User) TableName() string { return "users" }

func (u *User) LoadAssociations() []string { return nil }

func (u *User) GetID() any { return u.ID }
