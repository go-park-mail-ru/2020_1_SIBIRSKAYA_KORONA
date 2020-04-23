package models

//go:generate easyjson -all
type Checklist struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name" gorm:"not null" faker:"word"`
	Items []Item `json:"items",omitempty gorm:"foreignkey:clid"`
	Tid   uint   `json:"tid" gorm:"not null"`
}

//easyjson:json
type Checklists []Checklist

func (checkilist *Checklists) TableName() string {
	return "checklists"
}
