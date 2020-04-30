package attach

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/attach_usecase_mock.go
type UseCase interface {
	Create(attachModel *models.AttachedFile, attachFile *multipart.FileHeader) error
	Get(tid uint) (models.AttachedFiles, error)
	GetByID(tid uint, fid uint) (*models.AttachedFile, error)
	Delete(fid uint) error
}
