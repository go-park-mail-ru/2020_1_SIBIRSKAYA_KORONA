package models

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
)

type Column struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"title" gorm:"not null"`
	Pos   float64 `json:"position" gorm:"not null"`
	Tasks []Task  `json:"tasks,omitempty" gorm:"foreignkey:cid"`
	Bid   uint    `json:"-" gorm:"not null"`
}

func (col *Column) TableName() string {
	return "columns"
}

func CreateColumn(ctx echo.Context) *Column {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil
	}
	defer ctx.Request().Body.Close()
	column := new(Column)
	if json.Unmarshal(body, column) != nil {
		return nil
	}
	return column
}
