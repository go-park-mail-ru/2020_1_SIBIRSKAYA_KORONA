package models

type Board struct {
	ID      uint     `json:"bid" gorm:"primary_key"`
	Name    string   `json:"name"`
	Cols    []string `json:"cols"`
	Admins  []User   `json:"admins" gorm:"many2many:board_admins;"`
	Members []User   `json:"members" gorm:"many2many:board_members;"`
	// Tasks []Task
}

func (b *Board) TableName() string {
	return "boards"
}
