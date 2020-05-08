package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"
)

type NotificationUseCase struct {
	notificationRepo notification.Repository
}

func CreateUseCase(notificationRepo_ notification.Repository) notification.UseCase {
	return &NotificationUseCase{notificationRepo: notificationRepo_}
}

func (notificationUseCase *NotificationUseCase) Pop(uid uint) (models.Events, bool) {
	return notificationUseCase.notificationRepo.Pop(uid)
}
