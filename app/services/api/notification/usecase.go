package notification

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Pop(uid uint) (models.Events, bool)
}
