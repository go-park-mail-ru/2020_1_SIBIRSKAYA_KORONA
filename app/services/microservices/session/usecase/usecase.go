package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type SessionUseCase struct {
	sessionRepo session.Repository
}

func CreateUseCase(sessionRepo_ session.Repository) session.UseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo_,
	}
}

func (sessionUseCase *SessionUseCase) Create(ses models.Session) error {
	err := sessionUseCase.sessionRepo.Create(ses)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (sessionUseCase *SessionUseCase) Get(sid string) (uint, error) {
	id, err := sessionUseCase.sessionRepo.Get(sid)
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	return id, nil
}

func (sessionUseCase *SessionUseCase) Delete(sid string) error {
	err := sessionUseCase.sessionRepo.Delete(sid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
