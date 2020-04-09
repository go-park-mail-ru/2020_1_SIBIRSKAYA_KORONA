package column

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/column_repo_mock.go
type Repository interface {
	Create(column *models.Column) error
	Get(cid uint) (*models.Column, error)
	GetTasksByID(cid uint) ([]models.Task, error)
	Delete(cid uint) error
}
