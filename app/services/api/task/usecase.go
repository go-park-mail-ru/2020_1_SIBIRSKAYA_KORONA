package task

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/task_usecase_mock.go
type UseCase interface {
	Create(task *models.Task) error
	Get(cid uint, tid uint) (*models.Task, error)
	Update(newTask models.Task) error
	Delete(tid uint) error
	Assign(tid uint, uid uint) error
	Unassign(tid uint, uid uint) error

	//Comments -------------------
	CreateComment(comment *models.Comment) error
	GetComments(tid uint, uid uint) (models.Comments, error)
}
