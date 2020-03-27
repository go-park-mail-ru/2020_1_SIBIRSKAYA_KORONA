package models

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"time"
)

type Task struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Description string    `json:"description"`
	Members     []*User   `json:"members" gorm:"many2many:task_members"` // TODO(Alexandr): preload
	Order       float32   `json:"order"`	// TODO(Alexandr): auto_increment
	CreatedAt   time.Time `json:"createdAt"`

	//Labels []*Label

}

func (t *Task) TablaName() string {
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
