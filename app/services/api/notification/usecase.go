package notification

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/notification_usecase_mock.go
type UseCase interface {
	Create(event *models.Event) error
	GetAll(uid uint) (models.Events, bool)
	UpdateAll(uid uint) error
	DeleteAll(uid uint) error
}
