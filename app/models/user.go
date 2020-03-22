package models

import (
	"encoding/json"
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

func CreateUser(reqBody []byte) *User {
	usr := new(User)
	if json.Unmarshal(reqBody, usr) != nil {
		return nil
	}
	return usr
}
