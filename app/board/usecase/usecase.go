package usecase

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"

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

func (boardUseCase *BoardUseCase) GetAll(sid string) ([]models.Board, []models.Board, *cstmerr.UseError) {
	usr := boardUseCase.GetUser(sid)
	if usr == nil {
		return nil, nil, &cstmerr.UseError{Err: models.ErrSessionNotExist, Code: http.StatusUnauthorized}
	}

	adminsBoard, membersBoard, repoErr := boardUseCase.boardRepo.GetAll(usr)
	var responseCode int
	switch repoErr.Err {
	case models.ErrDbBadOperation:
		responseCode = http.StatusInternalServerError
		repoErr.Err = models.ErrInternal
	case nil:
		responseCode = http.StatusOK
	}

	return adminsBoard, membersBoard, &cstmerr.UseError{Err: repoErr.Err, Code: responseCode}
}
