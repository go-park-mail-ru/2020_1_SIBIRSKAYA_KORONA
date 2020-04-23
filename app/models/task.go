package models

//go:generate easyjson -all
type Task struct {
	ID         uint        `json:"id" gorm:"primary_key"`
	Name       string      `json:"title" gorm:"not null" faker:"word"`
	About      string      `json:"description" faker:"sentence"`
	Level      uint        `json:"level,omitempty"`
	Deadline   string      `json:"deadline,omitempty" faker:"date"`
	Pos        float64     `json:"position" gorm:"not null"`
	Cid        uint        `json:"cid" gorm:"not null"`
	Members    []User      `json:"members,omitempty" gorm:"many2many:task_members"`
	Comments   []Comment   `json:"comments,omitempty" gorm:"foreignkey:tid"`
	Checklists []Checklist `json:"checklists,omitempty" gorm:"foreignkey:tid"`
	Labels   []Label `json:"labels,omitempty" gorm:"many2many:task_labels;"`
}

//easyjson:json
type Tasks []Task

func (tsk *Task) TableName() string {
	return "tasks"
}
