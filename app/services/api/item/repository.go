package item

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(chlist *models.Item) error
	Update(chlist *models.Item) error
	Delete(itid uint) error
}
