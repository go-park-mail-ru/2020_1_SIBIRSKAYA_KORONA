package usecase

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"
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

func (sessionUseCase *SessionUseCase) Create(user *models.User, sessionExpires time.Time) (string, *cstmerr.UseError) {
	if user == nil {
		return "", &cstmerr.UseError{Err: models.ErrUserNotExist, Code: http.StatusUnauthorized}
	}

	realUser := sessionUseCase.userRepo.GetByNickname(user.Nickname)
	if realUser != nil && realUser.Password == user.Password {
		ses := &models.Session{
			SID:     "",
			ID:      realUser.ID,
			Expires: sessionExpires,
		}

		var responseStatus int
		sid, repoErr := sessionUseCase.sessionRepo.Create(ses)
		switch repoErr.Err {
		case models.ErrWrongPassword:
			responseStatus = http.StatusUnauthorized
		case models.ErrInternal:
			responseStatus = http.StatusInternalServerError
		case nil:
			responseStatus = http.StatusOK
		default:
			responseStatus = http.StatusInternalServerError
		}

		return sid, &cstmerr.UseError{Err: repoErr.Err, Code: responseStatus}
	}

	return "", &cstmerr.UseError{Err: models.ErrWrongPassword, Code: http.StatusUnauthorized}
}

func (sessionUseCase *SessionUseCase) Has(sid string) bool {
	_, has := sessionUseCase.sessionRepo.Get(sid)
	return has
}

func (sessionUseCase *SessionUseCase) Delete(sid string) error {
	return sessionUseCase.sessionRepo.Delete(sid)
}
