package models

//go:generate easyjson -all
type AttachedFile struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	URL     string `json:"url" gorm:"not null" faker:"word"`
	Name    string `json:"filename" gorm:"not null" faker:"word"`
	FileKey string `json:"-" gorm:"not null" faker:"word"`
	Tid     uint   `json:"-" gorm:"not null"`
}

//easyjson:json
type AttachedFiles []AttachedFile

func (AttachedFile *AttachedFile) TableName() string {
	return "attached_files"
}
