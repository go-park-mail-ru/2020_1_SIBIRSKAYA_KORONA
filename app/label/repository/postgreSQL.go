package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/label"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) label.Repository {
	return &UserStore{DB: db}
}
