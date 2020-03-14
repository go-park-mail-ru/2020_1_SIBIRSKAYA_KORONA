package session

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(user *models.User) (string, error)
	LogIn(id uint) (string, error)
	Has(sid string) bool
	LogOut(sid string) error
	Delete(id uint, sid string) error
}