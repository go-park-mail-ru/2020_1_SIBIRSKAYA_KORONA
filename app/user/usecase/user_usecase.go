package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type UserUseCase struct {
	sessionRepo session.Repository
	userRepo user.Repository
}

func CreateUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository) user.UseCase {
	return &UserUseCase{
		sessionRepo: sessionRepo_,
		userRepo: userRepo_,
	}
}

func (userUseCase *UserUseCase) Create(user *models.User) (string, error) {
	err := userUseCase.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	session := &models.Session{
		SID:     "",
		ID:      user.ID,
		Expires: time.Now().AddDate(1,0,0),
	}
	return userUseCase.sessionRepo.Create(session)
}

func (userUseCase *UserUseCase) Get(userKey string) *models.User {
	var id uint
	_, err := fmt.Sscan(userKey, &id)
	if err == nil {
		return userUseCase.userRepo.GetByID(id)
	}
	return userUseCase.userRepo.GetByNickName(userKey)
}

func (userUseCase *UserUseCase) GetAll(sid string) *models.User {
	if id, has := userUseCase.sessionRepo.Get(sid); has {
		return userUseCase.userRepo.GetAll(id)
	}
	return nil
}

func (userUseCase *UserUseCase) Update(sid string, oldPass string, newUser *models.User) error {
	id, has := userUseCase.sessionRepo.Get(sid)
	if !has {
		return errors.New("no user")
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