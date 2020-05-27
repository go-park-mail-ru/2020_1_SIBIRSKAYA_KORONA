package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/board_usecase_mock.go
type UseCase interface {
	Create(uid uint, board *models.Board) error
	GetBoardsByUser(uid uint) (models.Boards, models.Boards, error)
	Get(uid uint, bid uint, isAdmin bool) (*models.Board, error)
	GetLabelsByID(bid uint) (models.Labels, error)
	GetColumnsByID(bid uint) (models.Columns, error)
	Update(board *models.Board) error
	Delete(bid uint) error
	InviteMember(bid uint, uid uint) error
	DeleteMember(bid uint, uid uint) error
	GetUsersForInvite(bid uint, nicknamePart string, limit uint) (models.Users, error)
	InviteMemberByLink(uid uint, link string) (*models.Board, error)
	UpdateInviteLink(bid uint) error
}
