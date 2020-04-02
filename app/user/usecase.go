package user

import (
	"mime/multipart"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type UseCase interface {
	Create(user *models.User, sessionExpires time.Time) (string, error)
	GetByUserKey(userKey string) (*models.User, error) // userKey - id, nickname
	GetByCookie(sid string) (*models.User, error)
	Update(sid string, oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error
	Delete(sid string) error
}
