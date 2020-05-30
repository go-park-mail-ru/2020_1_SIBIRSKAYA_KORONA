package models

//go:generate easyjson -all
type Template struct {
	Variant string `json:"template" gorm:"not null" faker:"word"`
}

//easyjson:json
type Templates []Task

func (template *Template) TableName() string {
	return "templates"
}
