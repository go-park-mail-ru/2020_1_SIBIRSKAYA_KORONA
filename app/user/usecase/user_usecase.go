package usecase

import (
	"strconv"

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
	return userUseCase.sessionRepo.Create(user.ID)
}

func (userUseCase *UserUseCase) Get(userKey string) *models.User {
	id, err := strconv.Atoi(userKey)
	if err == nil {
		return userUseCase.userRepo.GetByID(uint(id))
	}
	return userUseCase.userRepo.GetByNickName(userKey)
}

func (userUseCase *UserUseCase) GetAll(sid string) *models.User {
	if id, has := userUseCase.sessionRepo.Get(sid); has {
		return userUseCase.userRepo.GetAll(id)
	}
	return nil
}

func (userUseCase *UserUseCase) Update(newUser *models.User) error {
	return userUseCase.userRepo.Update(newUser)
}

func (userUseCase *UserUseCase) Delete(id uint, sid string) error {
	err := userUseCase.sessionRepo.Delete(sid)
	if err != nil {
		return err
	}
	return userUseCase.userRepo.Delete(id)
}