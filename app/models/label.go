package models

//go:generate easyjson -all
type Label struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"title" gorm:"not null" faker:"word"`
	Color string `json:"color" faker:"sentence"`
	Bid   uint   `json:"-" gorm:"not null"`
	Tasks []Task `json:"-" gorm:"many2many:task_labels;" faker:"-"`
}

//easyjson:json
type Labels []Label

func (lbl *Label) TableName() string {
	return "labels"
}
