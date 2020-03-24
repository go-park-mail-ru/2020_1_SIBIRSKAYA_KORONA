package usecase

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/custom_errors"
)

// Создаем ошибки, характерные для юзкейса пользователя
var (
	ErrUserNotExist    = errors.New("User not exist")
	ErrUserBadMarshall = errors.New("Invalid data for user update")
)

type UserUseCase struct {
	sessionRepo session.Repository
	userRepo    user.Repository
}

func CreateUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository) user.UseCase {
	return &UserUseCase{
		sessionRepo: sessionRepo_,
		userRepo:    userRepo_,
	}
}

func (userUseCase *UserUseCase) Create(user *models.User, sessionExpires time.Time) (string, error) {
	err := userUseCase.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	ses := &models.Session{
		SID:     "",
		ID:      user.ID,
		Expires: sessionExpires,
	}
	return userUseCase.sessionRepo.Create(ses)
}

func (userUseCase *UserUseCase) GetByUserKey(userKey string) *models.User {
	var id uint
	_, err := fmt.Sscan(userKey, &id)
	usr := new(models.User)
	if err == nil {
		usr = userUseCase.userRepo.GetByID(id)
	} else {
		usr = userUseCase.userRepo.GetByNickname(userKey)
	}
	if usr != nil {
		usr.Password = ""
	}
	return usr
}

func (userUseCase *UserUseCase) GetByCookie(sid string) *models.User {
	id, has := userUseCase.sessionRepo.Get(sid)
	if !has {
		return nil
	}
	usr := userUseCase.userRepo.GetByID(id)
	if usr != nil {
		usr.Password = ""
	}
	return usr
}

func (userUseCase *UserUseCase) Update(sid string, oldPass string, newUser *models.User) error {
	if newUser == nil {
		return &custom_errors.CustomUsecaseError{Err: ErrUserBadMarshall, Code: http.StatusBadRequest}
	}
	id, has := userUseCase.sessionRepo.Get(sid)
	if !has {
		return &custom_errors.CustomUsecaseError{Err: ErrUserNotExist, Code: http.StatusBadRequest}
	}
	newUser.ID = id
	return userUseCase.userRepo.Update(oldPass, newUser)
}

func (userUseCase *UserUseCase) Delete(sid string) error {
	id, has := userUseCase.sessionRepo.Get(sid)
	if !has {
		return errors.New("no user")
	}
	err := userUseCase.sessionRepo.Delete(sid)
	if err != nil {
		return err
	}
	return userUseCase.userRepo.Delete(id)
}
