package column

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/column_usecase_mock.go
type UseCase interface {
	Create(column *models.Column) error
	Get(bid uint, cid uint) (*models.Column, error)
	GetTasksByID(cid uint) (models.Tasks, error)
	Update(newCol models.Column) error
	Delete(cid uint) error
}
