package task

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(task *models.Task) error
	Get(tid uint) (*models.Task, error)
}
