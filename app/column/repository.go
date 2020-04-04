package column

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(column *models.Column) error
	Get(cid uint) (*models.Column, error)
	GetTasksByID(cid uint) ([]models.Task, error)
}
