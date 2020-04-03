package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
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
	user, err := boardUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	board.Admins = []*models.User{user}
	return boardUseCase.boardRepo.Create(board)
}

func (boardUseCase *BoardUseCase) Update(uid uint, newBoard *models.Board) error {
	err := boardUseCase.boardRepo.Update(newBoard)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (boardUseCase *BoardUseCase) Delete(uid uint, bid uint) error {
	// _, err := boardUseCase.Get(uid, bid, true)
	// if err != nil {
	// 	logger.Error(err)
	// 	return err
	// }

	err := boardUseCase.boardRepo.Delete(bid)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (boardUseCase *BoardUseCase) Get(id uint, bid uint, isAdmin bool) (*models.Board, error) {
	brd, err := boardUseCase.boardRepo.Get(bid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	tmp := brd.Admins
	if !isAdmin {
		tmp = append(tmp, brd.Members...)
	}
	for _, member := range tmp {
		if member.ID == id {
			return brd, nil
		}
	}
	return nil, errors.ErrBoardsNotFound
}

func (boardUseCase *BoardUseCase) GetAll(uid uint) ([]models.Board, []models.Board, error) {
	user, err := boardUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	adminsBoard, membersBoard, repoErr := boardUseCase.boardRepo.GetAll(user)
	if repoErr != nil {
		logger.Error(repoErr)
		return nil, nil, repoErr
	}

	return adminsBoard, membersBoard, nil
}
