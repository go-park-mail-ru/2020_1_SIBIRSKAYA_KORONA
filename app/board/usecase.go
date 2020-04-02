package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type UseCase interface {
	Create(sid string, board *models.Board) error
	Get(sid string, bid uint) (*models.Board, error)
	GetAll(sid string) ([]models.Board, []models.Board, error)
	Update(sid string, board *models.Board) error
	Delete(sid string, bid uint) error
}
