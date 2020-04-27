package attach

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/board_usecase_mock.go
type Repository interface {
	Create(attach *models.AttachedFile) error
	Get(tid uint) (models.AttachedFiles, error)
	GetByID(fid uint) (*models.AttachedFile, error)
	Delete(fid uint) error
}
