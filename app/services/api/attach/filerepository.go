package attach

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=filerepository.go -package=mocks -destination=./mocks/attach_filerepo_mock.go
type FileRepository interface {
	UploadFile(attachFile *multipart.FileHeader, attachModel *models.AttachedFile) (string, error)
	DeleteFile(filenameKey string) error
}
