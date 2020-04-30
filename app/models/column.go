package models

//go:generate easyjson -all
type Column struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"title" gorm:"not null" faker:"word"`
	Pos   float64 `json:"position" gorm:"not null"`
	Tasks []Task  `json:"tasks,omitempty" gorm:"foreignkey:cid"`
	Bid   uint    `json:"-" gorm:"not null"`
}

//easyjson:json
type Columns []Column

func (col *Column) TableName() string {
	return "columns"
}
