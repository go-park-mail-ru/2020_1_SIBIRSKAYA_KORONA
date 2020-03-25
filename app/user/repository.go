package user

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"
)

type Repository interface {
	Create(user *models.User) error
	GetByID(id uint) *models.User
	GetByNickname(nickname string) *models.User
	Update(oldPass string, newUser *models.User) *cstmerr.CustomRepositoryError
	Delete(id uint) error
}
