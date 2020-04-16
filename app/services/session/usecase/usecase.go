package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session"
)

type SessionUseCase struct {
	sessionRepo session.Repository
}

func CreateUseCase(sessionRepo_ session.Repository) session.UseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo_,
	}
}

func (sessionUseCase *SessionUseCase) Create(sid string, uid uint32, expiration int32) error {
	return sessionUseCase.sessionRepo.Create(sid, uid, expiration)
}

func (sessionUseCase *SessionUseCase) Get(sid string) (uint, error) {
	return sessionUseCase.sessionRepo.Get(sid)
}

func (sessionUseCase *SessionUseCase) Delete(sid string) error {
	return sessionUseCase.sessionRepo.Delete(sid)
}
