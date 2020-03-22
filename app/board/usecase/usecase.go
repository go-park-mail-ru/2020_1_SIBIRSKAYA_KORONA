package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"

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

func (boardUseCase *BoardUseCase) GetUser(sid string) *models.User {
	id, has := boardUseCase.sessionRepo.Get(sid)
	if !has {
		return nil
	}
	return boardUseCase.userRepo.GetByID(id)
}
func (boardUseCase *BoardUseCase) Create(sid string, board *models.Board) error {
	usr := boardUseCase.GetUser(sid)
	if usr == nil {
		return errors.New("not found")
	}
	board.Admins = []models.User{*usr}
	return boardUseCase.boardRepo.Create(board)
}

func (boardUseCase *BoardUseCase) GetAll(sid string) ([]models.Board, []models.Board, error) {
	usr := boardUseCase.GetUser(sid)
	if usr == nil {
		return nil, nil, errors.New("not found")
	}
	return boardUseCase.boardRepo.GetAll(usr)
}
