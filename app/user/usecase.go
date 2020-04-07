package user

import (
	"mime/multipart"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/user_usecase_mock.go
type UseCase interface {
	Create(user *models.User, sessionExpires time.Time) (string, error)
	GetByID(uid uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	GetBoardsByID(uid uint) ([]models.Board, []models.Board, error)
	Update(oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error
	Delete(uid uint, sid string) error
}
