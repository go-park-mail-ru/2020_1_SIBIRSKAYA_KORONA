package task

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(sid string, task *models.Task) error
	Update(sid string, task *models.Task) error
	Delete(sid string, task *models.Task) error

	GetByBoardID(sid string, bid uint) []*models.Task
	GetByColumnID(sid string, cid uint) []*models.Task
	GetByID(sid string, tid uint) *models.Task
}