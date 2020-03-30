package board

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	Delete(bid uint) error

	Get(bid uint) (*models.Board, error)
	GetAll(user *models.User) ([]models.Board, []models.Board, error)
}
