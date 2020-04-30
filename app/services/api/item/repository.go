package item

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/item_repo_mock.go
type Repository interface {
	Create(chlist *models.Item) error
	Update(chlist *models.Item) error
	Delete(itid uint) error
	GetByID(itid uint) (*models.Item, error)
}
