package models

type User struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	PathToAvatar string `json:"avatar"`
	Password     string `json:"password,omitempty"`
}
