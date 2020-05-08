package models

//go:generate easyjson -all
type Event struct {
	ID       uint    `json:"-" gorm:"primary_key"`
	Message  string  `json:"message" gorm:"not null"`
}

//easyjson:json
type Events []Event

func (event *Event) TableName() string {
	return "events"
}
