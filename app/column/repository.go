package column

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(column *models.Column) error
	Update(column *models.Column) error
	Delete(cid uint) error

	GetByBoardID(bid uint) ([]*models.Column, error)
	GetByID(cid uint) (*models.Column, error)
}