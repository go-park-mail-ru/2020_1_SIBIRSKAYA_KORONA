package models

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"time"
)

type Column struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Title     *string    `json:"title"`
	Tasks     []*Task   `json:"tasks"`
	Position     *float32   `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
}

func (c *Column) TableName() string {
	return "columns"
}

func CreateColumn(ctx echo.Context) *Column {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err!= nil {
		return nil
	}
	defer ctx.Request().Body.Close()

	column := new(Column)
	if json.Unmarshal(body, column) != nil {
		return nil
	}
	return column
}
