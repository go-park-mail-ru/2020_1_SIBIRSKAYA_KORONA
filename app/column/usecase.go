package column

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(column *models.Column) error
	Get(bid uint, cid uint) (*models.Column, error)
	GetTasksByID(cid uint) ([]models.Task, error)
}
