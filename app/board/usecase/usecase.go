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
	board.Admins = []*models.User{usr}
	return boardUseCase.boardRepo.Create(board)
}

func (boardUseCase *BoardUseCase) Update(sid string, newBoard *models.Board) error {
	usr := boardUseCase.GetUser(sid)
	if usr == nil {
		return errors.New("not found")
	}

	oldBoard, err := boardUseCase.boardRepo.Get(newBoard.ID)
	isAdmin := false

	for _, admin := range oldBoard.Admins {
		if usr.ID == admin.ID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return errors.New("no permission")
	}

	return boardUseCase.boardRepo.Update(newBoard)
}

func (boardUseCase *BoardUseCase) Delete(sid string, bid uint) error {
	usr := boardUseCase.GetUser(sid)
	if usr == nil {
		return errors.New("not found")
	}

	boardToDelete := boardUseCase.boardRepo.Get(bid)
	isAdmin := false

	for _, admin := range boardToDelete.Admins {
		if usr.ID == admin.ID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return errors.New("no permission to delete")
	}

	return boardUseCase.boardRepo.Delete(bid)
}

func (boardUseCase *BoardUseCase) Get(sid string, bid uint) *models.Board {
	usr := boardUseCase.GetUser(sid)
	if usr == nil {
		return nil
	}
	brd := boardUseCase.boardRepo.Get(bid)
	for _, member := range append(brd.Admins, brd.Members...) {
		if member.ID == usr.ID {
			return brd
		}
	}
	return nil
}

func (boardUseCase *BoardUseCase) GetAll(sid string) ([]models.Board, []models.Board, error) {
	usr := boardUseCase.GetUser(sid)
	if usr == nil {
		return nil, nil, errors.New("not found")
	}
	return boardUseCase.boardRepo.GetAll(usr)
}
