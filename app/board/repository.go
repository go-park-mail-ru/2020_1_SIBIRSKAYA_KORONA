package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type Repository interface {
	Create(board *models.Board) error
	Get(bid uint) (*models.Board, error)
	GetColumnsByID(bid uint) ([]models.Column, error)
	Update(board *models.Board) error
	Delete(bid uint) error
}
