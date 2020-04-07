package task

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/task_repo_mock.go
type Repository interface {
	Create(task *models.Task) error
	Get(tid uint) (*models.Task, error)
	Update(newTask models.Task) error
	Delete(tid uint) error
}
