package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID       uint     `json:"id" gorm:"primary_key"`
	Name     string   `json:"name" gorm:"not null" faker:"name"`
	Surname  string   `json:"surname" gorm:"not null" faker:"last_name"`
	Nickname string   `json:"nickname" gorm:"unique;not null" faker:"username"`
	Email    string   `json:"email" faker:"email" faker:"email"`
	Avatar   string   `json:"avatar" faker:"url"`
	Password string   `json:"password,omitempty" gorm:"not null" faker:"password"`
	Admin    []*Board `json:"-" gorm:"many2many:board_admins;" faker:"-"`
	Member   []*Board `json:"-" gorm:"many2many:board_members;" faker:"-"`
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
