package user

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Get(id uint) *models.User // only public info
	GetAll(id uint) *models.User
	Update(newUser *models.User) error
}
