package user

import (
	"mime/multipart"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"
)

type UseCase interface {
	Create(user *models.User, sessionExpires time.Time) (string, error)
	GetByUserKey(userKey string) *models.User // userKey - id, nickname
	GetByCookie(sid string) *models.User
	Update(sid string, oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) *cstmerr.UseError
	Delete(sid string) error
}
