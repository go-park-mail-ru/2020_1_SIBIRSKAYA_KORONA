package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/board_usecase_mock.go
type UseCase interface {
	Create(uid uint, board *models.Board) error
	Get(uid uint, bid uint, isAdmin bool) (*models.Board, error)
	GetColumnsByID(bid uint) ([]models.Column, error)
	Update(board *models.Board) error
	Delete(bid uint) error
	InviteMember(bid uint, uid uint) error
	DeleteMember(bid uint, uid uint) error
}
