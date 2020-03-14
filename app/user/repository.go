package user

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(user *models.User) error
	Get(id uint) *models.User // only public info
	GetAll(id uint) *models.User
	Update(newUser *models.User) error
	Delete(id uint) error
}
