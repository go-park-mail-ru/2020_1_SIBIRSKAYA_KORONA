package models

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
)

//go:generate easyjson -all
type User struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Name     string  `json:"name" gorm:"not null" faker:"name"`
	Surname  string  `json:"surname" gorm:"not null" faker:"last_name"`
	Nickname string  `json:"nickname" gorm:"unique;not null" faker:"username"`
	Avatar   string  `json:"avatar" faker:"url"`
	Email    string  `json:"email,omitempty" faker:"email"`
	Password []byte  `json:"password,omitempty" gorm:"not null" faker:"-"`
	Admin    []Board `json:"-" gorm:"many2many:board_admins;" faker:"-"`
	Member   []Board `json:"-" gorm:"many2many:board_members;" faker:"-"`
}

//easyjson:json
type Users []User

func (usr *User) TableName() string {
	return "users"
}

// TODO: не маршалить

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
	return &res
}
