package task

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(task *models.Task) error
	Get(cid uint, tid uint) (*models.Task, error)
	Update(newTask models.Task) error
	Delete(tid uint) error
}
