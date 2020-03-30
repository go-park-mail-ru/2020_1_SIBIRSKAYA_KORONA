package column

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(sid string, bid uint, column *models.Column) error
	Update(sid string, column *models.Column) error
	Delete(sid string, cid uint) error

	GetByBoardID(sid string, bid uint) ([]*models.Column, error)
	GetByID(sid string, cid uint) (*models.Column, error)
}
