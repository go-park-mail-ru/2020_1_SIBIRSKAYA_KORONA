package models

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/sanitize"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Name     string  `json:"name" gorm:"not null" faker:"name"`
	Surname  string  `json:"surname" gorm:"not null" faker:"last_name"`
	Nickname string  `json:"nickname" gorm:"unique;not null" faker:"username"`
	Email    string  `json:"email,omitempty" faker:"email"`
	Avatar   string  `json:"avatar" faker:"url"`
	Password []byte  `json:"password,omitempty" gorm:"not null" faker:"-"`
	Admin    []Board `json:"-" gorm:"many2many:board_admins;" faker:"-"`
	Member   []Board `json:"-" gorm:"many2many:board_members;" faker:"-"`
}

//easyjson:json
type Users []User

func (usr *User) TableName() string {
	return "users"
}

func (usr *User) ToProto() *proto.UserMess {
	usrJson, err := json.Marshal(usr)
	if err != nil {
		return nil
	}
	var res proto.UserMess
	err = json.Unmarshal(usrJson, &res)
	if err != nil {
		return nil
	}
	res.Password = usr.Password
	return &res
}

func CreateUserFromProto(usr proto.UserMess) *User {
	usrJson, err := json.Marshal(usr)
	if err != nil {
		return nil
	}
	var res User
	err = json.Unmarshal(usrJson, &res)
	if err != nil {
		return nil
	}
	res.Password = usr.Password
	return &res
}

func CreateUser(ctx echo.Context) *User {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil
	}
	defer ctx.Request().Body.Close()
	sanBody, err := sanitize.SanitizeJSON(body)
	if err != nil {
		return nil
	}
	usr := new(User)
	if json.Unmarshal(sanBody, usr) != nil {
		return nil
	}
	return usr
}
