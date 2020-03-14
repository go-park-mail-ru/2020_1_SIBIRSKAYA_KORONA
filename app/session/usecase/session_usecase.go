package usecase

import (
	//"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
)

type sessionUsecase struct {
	sessionRepo_ session.Repository
}

func NewSessionUsecase(sessionRepo session.Repository) session.Usecase {
	return &sessionUsecase{sessionRepo_: sessionRepo}
}

func (this *sessionUsecase) AddSession(nickname string) string {
	return this.sessionRepo_.AddSession(string)
}

func (this *sessionUsecase) GetSession(SID string) (string, bool) {
	return this.sessionRepo_.GetSession(SID)
}

func (this *sessionUsecase) DeleteSession(SID string) error {
	return this.sessionRepo_.DeleteSession(SID)
}
