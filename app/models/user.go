package models

import "encoding/json"

type User struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	Img string `json:"avatar"`
	Password     string `json:"password,omitempty"`
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

