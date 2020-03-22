package user

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(user *models.User) error
	GetByID(id uint) *models.User
	GetByNickname(nickname string) *models.User
	Update(oldPass string, newUser *models.User) error
	Delete(id uint) error
}
