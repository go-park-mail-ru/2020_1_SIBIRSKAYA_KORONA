package models

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
)

type Task struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Name     string  `json:"title" gorm:"not null"`
	About    string  `json:"description"`
	Level    uint    `json:"level,omitempty"`
	Deadline string  `json:"deadline,omitempty"`
	Pos      float64 `json:"position" gorm:"not null"`
	Cid      uint    `json:"cid" gorm:"not null"`
	// Members     []User   `json:"members,omitempty" gorm:"many2many:task_members"`
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
	task := new(Task)
	if json.Unmarshal(body, task) != nil {
		return nil
	}
	return task
}
