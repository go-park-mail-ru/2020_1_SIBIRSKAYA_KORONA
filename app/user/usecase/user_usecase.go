package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type userUsecase struct {
	userRepo_ user.Repository 
}

func NewUserUsecase(userRepo user.Repository) user.Usecase {
	return &userUsecase{userRepo_: userRepo}
}

func (this *userUsecase) AddUser(user *models.User) {
	this.userRepo_.AddUser(user)
}

func (this *userUsecase) GetUser(nickname string) (*models.User, bool) {
	return this.userRepo_.GetUser(nickname)
}