package user

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type Repository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	Update(oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error
	Delete(id uint) error
}
