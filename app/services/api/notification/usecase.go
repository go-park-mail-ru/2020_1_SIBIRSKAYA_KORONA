package notification

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(event *models.Event) error
	GetAll(uid uint) (models.Events, bool)
	UpdateAll(uid uint) error
	DeleteAll(uid uint) error
}
