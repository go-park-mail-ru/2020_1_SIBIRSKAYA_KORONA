package models

//go:generate easyjson -all
type Comment struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	Text           string `json:"text" gorm:"not null" faker:"sentence"`
	CreatedAt      int64  `json:"createdAt,omitempty" gorm:"not null" faker:"date"`
	IsEdited       bool   `json:"edited" gorm:"not null"`
	Uid            uint   `json:"-" gorm:"not null"`
	Nickname       string `json:"nickname" gorm:"-"`
	Avatar         string `json:"avatar" gorm:"-"`
	Tid            uint   `json:"-" gorm:"not null"`
	ReaderIsAuthor bool   `json:"readerIsAuthor" gorm:"-"`
}

//easyjson:json
type Comments []Comment

func (com *Comment) TableName() string {
	return "comments"
}
