package task

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(task *models.Task) error
	Update(task *models.Task) error
	Delete(task *models.Task) error

	GetByBoardID(bid uint) ([]*models.Task, error)
	GetByColumnID(cid uint) ([]*models.Task, error)
	GetByID(tid uint) (*models.Task, error)
}
