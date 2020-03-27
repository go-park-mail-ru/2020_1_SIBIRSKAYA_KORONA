package column

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(sid string, column *models.Column) error
	Update(sid string, column *models.Column) error
	Delete(sid string, column *models.Column) error

	GetByBoardID(sid string, bid uint) []*models.Column
	GetByID(sid string, cid uint) *models.Column
}
