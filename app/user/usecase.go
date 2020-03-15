package user

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(user *models.User) (string, error)
	Get(userKey string) *models.User // userKey - id, nickname
	GetAll(sid string) *models.User
	Update(newUser *models.User) error
	Delete(id uint, sid string) error
}
