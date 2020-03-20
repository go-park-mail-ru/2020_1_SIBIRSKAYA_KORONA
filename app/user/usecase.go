package user

import (
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type UseCase interface {
	Create(user *models.User, sessionExpires time.Time) (string, error)
	Get(userKey string) *models.User // userKey - id, nickname
	GetAll(sid string) *models.User
	Update(sid string, oldPass string, newUser *models.User) error
	Delete(sid string) error
}
