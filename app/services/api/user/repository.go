package user

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"mime/multipart"
)

type Repository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	CheckPassword(uid uint, pass []byte) bool
	Update(oldPass []byte, newUser models.User, avatarFileDescriptor *multipart.FileHeader) error
	Delete(id uint) error
	GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error)
}
