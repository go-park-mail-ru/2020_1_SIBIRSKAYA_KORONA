package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

type Board struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	// Cols    []string `json:"cols"`
	Admins  []*User `json:"admins" gorm:"many2many:board_admins;"`
	Members []*User `json:"members" gorm:"many2many:board_members;"`
	// Members []*User `json:"members,omitempty" gorm:"many2many:board_members;"`
	// Tasks []Task
}

func (b *Board) TableName() string {
	return "boards"
}

func CreateBoard(ctx echo.Context) *Board {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil
	}
	defer ctx.Request().Body.Close()
	board := new(Board)
	if json.Unmarshal(body, board) != nil {
		return nil
	}
	return board
}
