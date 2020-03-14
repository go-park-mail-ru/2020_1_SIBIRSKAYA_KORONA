package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type SessionUseCase struct {
	sessionRepo session.Repository
	userRepo user.Repository
}

func CreateSessionUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository) session.UseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo_,
		userRepo: userRepo_,
	}
}

func (sessionUseCase *SessionUseCase) Create(user *models.User) (string, error) {
	err := sessionUseCase.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	return sessionUseCase.sessionRepo.Create(user.ID)
}

func (sessionUseCase *SessionUseCase) LogIn(id uint) (string, error) {
	return sessionUseCase.sessionRepo.Create(id)
}

func (sessionUseCase *SessionUseCase) Has(sid string) bool {
	return sessionUseCase.sessionRepo.Has(sid)
}

func (sessionUseCase *SessionUseCase) LogOut(sid string) error {
	return sessionUseCase.sessionRepo.Delete(sid)
}

func (sessionUseCase *SessionUseCase) Delete(id uint, sid string) error {
	err := sessionUseCase.sessionRepo.Delete(sid)
	if err != nil {
		return err
	}
	return sessionUseCase.userRepo.Delete(id)
}