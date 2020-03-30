package task

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(sid string, task *models.Task) error
	Update(sid string, task *models.Task) error
	Delete(sid string, tid uint) error

	GetByBoardID(sid string, bid uint) ([]*models.Task, error)
	GetByColumnID(sid string, cid uint) ([]*models.Task, error)
	GetByID(sid string, tid uint) (*models.Task, error)
}
