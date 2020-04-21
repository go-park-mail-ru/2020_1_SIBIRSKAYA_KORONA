package models

//go:generate easyjson -all
type Board struct {
	ID      uint     `json:"id" gorm:"primary_key"`
	Name    string   `json:"title" gorm:"not null" faker:"word"`
	Columns []Column `json:"columns,omitempty" gorm:"foreignkey:bid" faker:"-"`
	Labels  []Label  `json:"labels,omitempty" gorm:"foreignkey:bid" faker:"-"`
	Admins  []User   `json:"admins,omitempty" gorm:"many2many:board_admins;" faker:"-"`
	Members []User   `json:"members,omitempty" gorm:"many2many:board_members;" faker:"-"`
}

//easyjson:json
type Boards []Board

type UserBoards struct {
	Admin  Boards `json:"admin"`
	Member Boards `json:"member"`
}

func (brd *Board) TableName() string {
	return "boards"
}
