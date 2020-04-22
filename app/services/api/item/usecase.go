package item

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(item *models.Item) error
	Update(item *models.Item) error
	Delete(itid uint) error
}
