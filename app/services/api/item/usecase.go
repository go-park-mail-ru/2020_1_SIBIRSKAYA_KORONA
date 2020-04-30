package item

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/item_usecase_mock.go
type UseCase interface {
	Create(itm *models.Item) error
	Update(itm *models.Item) error
	Delete(itid uint) error
	GetByID(clid uint, itid uint) (*models.Item, error)
}
