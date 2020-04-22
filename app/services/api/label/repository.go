package label

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/user_repo_mock.go
type Repository interface {
	Create(lbl *models.Label) error
	Get(lid uint) (*models.Label, error)
	Update(lbl models.Label) error
	Delete(lid uint) error
}