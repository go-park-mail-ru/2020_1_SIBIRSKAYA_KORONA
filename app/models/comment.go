package models

//go:generate easyjson -all
type Comment struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Text      string `json:"text" gorm:"not null" faker:"sentence"`
	CreatedAt int64  `json:"createdAt,omitempty" gorm:"not null" faker:"date"`
	IsEdited  bool   `json:"edited" gorm:"not null"`
	Uid       uint   `json:"-" gorm:"not null"`
	Tid       uint   `json:"-" gorm:"not null"`
}

//easyjson:json
type Comments []Comment

func (com *Comment) TableName() string {
	return "comments"
}

// Text
// CreatedAt
// IsEdited
// uid (not null)
// tid (not null)
