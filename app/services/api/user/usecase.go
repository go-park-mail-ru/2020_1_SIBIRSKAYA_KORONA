package user

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"mime/multipart"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/user_usecase_mock.go
type UseCase interface {
	Create(user *models.User, sessionExpires int32) (string, error)
	GetByID(uid uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error)
	Update(oldPass []byte, newUser models.User, avatarFileDescriptor *multipart.FileHeader) error
	Delete(uid uint, sid string) error
}
