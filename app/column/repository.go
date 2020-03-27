package column

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(column *models.Column) error
	Update(column *models.Column) error
	Delete(column *models.Column) error

	GetByBoardID(bid uint) []*models.Column
	GetByID(cid uint) *models.Column
}