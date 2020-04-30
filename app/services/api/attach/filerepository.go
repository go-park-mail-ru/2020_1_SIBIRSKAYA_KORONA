package attach

import (
	"mime/multipart"
)

//go:generate mockgen -source=filerepository.go -package=mocks -destination=./mocks/attach_filerepo_mock.go
type FileRepository interface {
	UploadFile(attachFile *multipart.FileHeader) (string, error)
	DeleteFile(filenameKey string) error
}
