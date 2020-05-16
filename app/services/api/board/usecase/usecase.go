package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type BoardUseCase struct {
	userRepo  user.Repository
	boardRepo board.Repository
}

func CreateUseCase(userRepo_ user.Repository, boardRepo_ board.Repository) board.UseCase {
	return &BoardUseCase{
		userRepo:  userRepo_,
		boardRepo: boardRepo_,
	}
}

func (boardUseCase *BoardUseCase) Create(uid uint, board *models.Board) error {
	usr, err := boardUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = boardUseCase.boardRepo.Create(uid, board)
	if err != nil {
		logger.Error(err)
		return err
	}
	board.Admins = append(board.Admins, *usr)
	return nil
}

func (boardUseCase *BoardUseCase) GetBoardsByUser(uid uint) (models.Boards, models.Boards, error) {
	adminsBoard, membersBoard, err := boardUseCase.boardRepo.GetBoardsByUser(uid)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}
	return adminsBoard, membersBoard, nil
}

func (boardUseCase *BoardUseCase) Get(id uint, bid uint, isAdmin bool) (*models.Board, error) {
	brd, err := boardUseCase.boardRepo.Get(bid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	usrs := brd.Admins
	if !isAdmin {
		usrs = append(usrs, brd.Members...)
	}
	for _, member := range usrs {
		if member.ID == id {
			return brd, nil
		}
	}
	return nil, errors.ErrNoPermission
}

func (boardUseCase *BoardUseCase) GetLabelsByID(bid uint) (models.Labels, error) {
	lbls, repoErr := boardUseCase.boardRepo.GetLabelsByID(bid)
	if repoErr != nil {
		logger.Error(repoErr)
		return nil, repoErr
	}
	return lbls, nil
}

func (boardUseCase *BoardUseCase) GetColumnsByID(bid uint) (models.Columns, error) {
	cols, repoErr := boardUseCase.boardRepo.GetColumnsByID(bid)
	if repoErr != nil {
		logger.Error(repoErr)
		return nil, repoErr
	}
	return cols, nil
}

func (boardUseCase *BoardUseCase) Update(newBoard *models.Board) error {
	err := boardUseCase.boardRepo.Update(newBoard)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (boardUseCase *BoardUseCase) Delete(bid uint) error {
	err := boardUseCase.boardRepo.Delete(bid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (boardUseCase *BoardUseCase) InviteMember(bid uint, uid uint) error {
	usr, err := boardUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = boardUseCase.boardRepo.InviteMember(bid, usr)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (boardUseCase *BoardUseCase) DeleteMember(bid uint, uid uint) error {
	usr, err := boardUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = boardUseCase.boardRepo.DeleteMember(bid, usr)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (boardUseCase *BoardUseCase) GetUsersForInvite(bid uint, nicknamePart string, limit uint) (models.Users, error) {
	users, err := boardUseCase.boardRepo.GetUsersForInvite(bid, nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (boardUseCase *BoardUseCase) InviteMemberByLink(uid uint, link string) (*models.Board, error) {
	usr, err := boardUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	brd, err := boardUseCase.boardRepo.InviteMemberByLink(*usr, link)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return brd, nil
}

func (boardUseCase *BoardUseCase) UpdateInviteLink(bid uint) error {
	err := boardUseCase.boardRepo.UpdateInviteLink(bid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
