package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type SessionUseCase struct {
	sessionRepo session.Repository
	userRepo    user.Repository
}

func CreateUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository) session.UseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo_,
		userRepo:    userRepo_,
	}
}

func (sessionUseCase *SessionUseCase) Create(user *models.User, sessionExpires time.Time) (string, error) {
	if user == nil {
		return "", errors.New("bad password")
	}
	realUser := sessionUseCase.userRepo.GetByNickname(user.Nickname)
	fmt.Println("!!!! pass !!!!", realUser.Password, user.Password)
	if realUser != nil && realUser.Password == user.Password {
		ses := &models.Session{
			SID:     "",
			ID:      realUser.ID,
			Expires: sessionExpires,
		}
		return sessionUseCase.sessionRepo.Create(ses)
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
