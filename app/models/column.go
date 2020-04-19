package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/sanitize"
	"github.com/labstack/echo/v4"
)

//go:generate easyjson -all
type Column struct {
	ID    uint    `json:"id" gorm:"primary_key"`
	Name  string  `json:"title" gorm:"not null" faker:"word"`
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

	sanBody, err := sanitize.SanitizeJSON(body)
	if err != nil {
		return nil
	}

	column := new(Column)
	if json.Unmarshal(sanBody, column) != nil {
		return nil
	}
	return column
}
