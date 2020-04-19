package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/board_repo_mock.go
type Repository interface {
	Create(uid uint, board *models.Board) error
	GetBoardsByUser(uid uint) ([]models.Board, []models.Board, error)
	Get(bid uint) (*models.Board, error)
	GetColumnsByID(bid uint) ([]models.Column, error)
	Update(board *models.Board) error
	Delete(bid uint) error
	InviteMember(bid uint, member *models.User) error
	DeleteMember(bid uint, member *models.User) error
}
