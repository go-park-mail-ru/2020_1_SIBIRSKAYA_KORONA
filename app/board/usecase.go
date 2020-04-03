package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type UseCase interface {
	Create(uid uint, board *models.Board) error
	Get(uid uint, bid uint, isAdmin bool) (*models.Board, error)
	GetAll(uid uint) ([]models.Board, []models.Board, error)
	Update(uid uint, board *models.Board) error
	Delete(uid uint, bid uint) error
}
