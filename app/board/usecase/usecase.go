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
	boardRepo board.Repository
}

func CreateUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository, boardRepo_ board.Repository) board.UseCase {
	return &BoardUseCase{
		sessionRepo: sessionRepo_,
		userRepo:    userRepo_,
		boardRepo:    boardRepo_,
	}
}

func (boardUseCase *BoardUseCase) Create(sid string, board *models.Board) error {
	id, has := boardUseCase.sessionRepo.Get(sid)
	if !has {
		return errors.New("no session")
	}
	user := boardUseCase.userRepo.GetByID(id)
	if user == nil {
		return errors.New("not found")
	}
	board.Admins = []models.User{*user}
	return boardUseCase.boardRepo.Create(board)
}
