package attach

import (
	"mime/multipart"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/board_usecase_mock.go
type FileRepository interface {
	UploadFile(attachFile *multipart.FileHeader) (string, error)
	DeleteFile(filenameKey string) error
}
