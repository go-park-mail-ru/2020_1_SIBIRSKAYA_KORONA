package models

//go:generate easyjson -all
type Item struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Text   string `json:"text" gorm:"not null" faker:"word"`
	IsDone bool   `json:"done" faker:"sentence"`
	Clid   uint   `json:"clid" gorm:"not null"`
}

//easyjson:json
type Items []Item

func (item *Item) TableName() string {
	return "items"
}
