package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/sanitize"
	"github.com/labstack/echo/v4"
)

type Task struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Name     string  `json:"title" gorm:"not null" faker:"word"`
	About    string  `json:"description" faker:"sentence"`
	Level    uint    `json:"level,omitempty"`
	Deadline string  `json:"deadline,omitempty" faker:"date"`
	Pos      float64 `json:"position" gorm:"not null"`
	Cid      uint    `json:"cid" gorm:"not null"`
	Members  []User  `json:"members,omitempty" gorm:"many2many:task_members"`
	// Labels []*Label
}

func (tsk *Task) TablaName() string {
	return "tasks"
}

func CreateTask(ctx echo.Context) *Task {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil
	}
	defer ctx.Request().Body.Close()

	sanBody, err := sanitize.SanitizeJSON(body)
	if err != nil {
		return nil
	}

	task := new(Task)
	if json.Unmarshal(sanBody, task) != nil {
		return nil
	}
	return task
}
