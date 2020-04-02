package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type BoardUseCase struct {
	sessionRepo session.Repository
	userRepo    user.Repository
	boardRepo   board.Repository
}

func CreateUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository, boardRepo_ board.Repository) board.UseCase {
	return &BoardUseCase{
		sessionRepo: sessionRepo_,
		userRepo:    userRepo_,
		boardRepo:   boardRepo_,
	}
}

func (boardUseCase *BoardUseCase) GetUser(sid string) (*models.User, error) {
	id, has := boardUseCase.sessionRepo.Get(sid)
	if !has {
		return nil, errors.ErrSessionNotExist
	}

	user, err := boardUseCase.userRepo.GetByID(id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (boardUseCase *BoardUseCase) Create(sid string, board *models.Board) error {
	user, err := boardUseCase.GetUser(sid)
	if err != nil {
		logger.Error(err)
		return err
	}
	board.Admins = []*models.User{user}
	return boardUseCase.boardRepo.Create(board)
}

func (boardUseCase *BoardUseCase) Get(sid string, bid uint) (*models.Board, error) {
	user, err := boardUseCase.GetUser(sid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	brd, err := boardUseCase.boardRepo.Get(bid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	for _, member := range append(brd.Admins, brd.Members...) {
		if member.ID == user.ID {
			return brd, nil
		}
	}
	return nil, errors.ErrBoardsNotFound
}

func (boardUseCase *BoardUseCase) GetAll(sid string) ([]models.Board, []models.Board, error) {
	user, err := boardUseCase.GetUser(sid)
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
