package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type UserUseCase struct {
	userRepo user.Repository
}

func CreateUserUseCase(userRepo_ user.Repository) user.UseCase {
	return &UserUseCase{userRepo: userRepo_}
}

func (userUseCase *UserUseCase) Get(id uint) *models.User {
	return userUseCase.userRepo.Get(id)
}

func (userUseCase *UserUseCase) GetAll(id uint) *models.User {
	return userUseCase.userRepo.GetAll(id)
}

func (userUseCase *UserUseCase) Update(newUser *models.User) error {
	return userUseCase.userRepo.Update(newUser)
}