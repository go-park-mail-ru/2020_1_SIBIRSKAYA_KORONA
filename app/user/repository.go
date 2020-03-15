package user

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(user *models.User) error
	GetByID(id uint) *models.User
	GetByNickName(nickName string) *models.User
	GetAll(id uint) *models.User
	Update(newUser *models.User) error
	Delete(id uint) error
}
