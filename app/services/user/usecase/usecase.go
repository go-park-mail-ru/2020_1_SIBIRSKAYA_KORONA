package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type UserUseCase struct {
	userRepo user.Repository
}

func CreateUseCase(userRepo_ user.Repository) user.UseCase {
	return &UserUseCase{
		userRepo: userRepo_,
	}
}

func (userUseCase *UserUseCase) Create(user *models.User) error {
	err := userUseCase.userRepo.Create(user)
	if err != nil {
		logger.Error(err)
	}
	return nil
}

func (userUseCase *UserUseCase) GetByID(uid uint) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (userUseCase *UserUseCase) GetByNickname(nickname string) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByNickname(nickname)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (userUseCase *UserUseCase) CheckPassword(uid uint, pass []byte) bool {
	return userUseCase.userRepo.CheckPassword(uid, pass)
}

func (userUseCase *UserUseCase) Update(oldPass []byte, newUser models.User) error {
	err := userUseCase.userRepo.Update(oldPass, newUser)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (userUseCase *UserUseCase) Delete(uid uint) error {
	err := userUseCase.userRepo.Delete(uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (userUseCase *UserUseCase) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	usrs, err := userUseCase.userRepo.GetUsersByNicknamePart(nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usrs, nil
}
