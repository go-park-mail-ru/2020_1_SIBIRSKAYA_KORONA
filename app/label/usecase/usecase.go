package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/label"
)

type UserUseCase struct {
	labelRepo label.Repository
}

func CreateUseCase(labelRepo_ label.Repository) label.UseCase {
	return &UserUseCase{
		labelRepo: labelRepo_,
	}
}
