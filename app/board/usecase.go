package board

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(sid string, board *models.Board) error
	GetAll(sid string) ([]models.Board, []models.Board, error)
}
