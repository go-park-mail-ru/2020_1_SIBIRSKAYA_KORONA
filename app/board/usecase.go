package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type UseCase interface {
	Create(uid uint, board *models.Board) error
	Get(uid uint, bid uint, isAdmin bool) (*models.Board, error)
	GetColumnsByID(bid uint) ([]models.Column, error)
	Update(board *models.Board) error
	Delete(bid uint) error
}
