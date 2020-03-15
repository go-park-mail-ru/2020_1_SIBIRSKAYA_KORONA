package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type SessionUseCase struct {
	sessionRepo session.Repository
	userRepo user.Repository
}

func CreateUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository) session.UseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo_,
		userRepo: userRepo_,
	}
}

func (sessionUseCase *SessionUseCase) Create(user *models.User) (string, error) {
	realUser := sessionUseCase.userRepo.GetByNickName(user.Nickname)
	if realUser.Password == user.Password {
		return sessionUseCase.sessionRepo.Create(realUser.ID)
	}
	return "", errors.New("bad password")
}

func (sessionUseCase *SessionUseCase) Has(sid string) bool {
	_, has := sessionUseCase.sessionRepo.Get(sid)
	return has
}

func (sessionUseCase *SessionUseCase) Delete(sid string) error {
	return sessionUseCase.sessionRepo.Delete(sid)
}