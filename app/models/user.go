package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Surname  string `json:"surname" gorm:"not null"`
	Nickname string `json:"nickname" gorm:"unique;not null"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Password string `json:"password,omitempty" gorm:"not null"`
}

func (u *User) TableName() string {
	return "users"
}

func CreateUser(ctx echo.Context) *User {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil
	}
	defer ctx.Request().Body.Close()
	usr := new(User)
	if json.Unmarshal(body, usr) != nil {
		return nil
	}
	return usr
}
