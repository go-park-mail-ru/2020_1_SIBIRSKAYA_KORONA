package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type NotificationUseCase struct {
	notificationRepo notification.Repository
	userRepo         user.Repository
}

func CreateUseCase(userRepo_ user.Repository, notificationRepo_ notification.Repository) notification.UseCase {
	return &NotificationUseCase{notificationRepo: notificationRepo_, userRepo: userRepo_}
}

func (notificationUseCase *NotificationUseCase) Create(event *models.Event) error {
	if event == nil {
		return errors.ErrInternal
	}
	if err := notificationUseCase.notificationRepo.Create(event); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (notificationUseCase *NotificationUseCase) GetAll(uid uint) (models.Events, bool) {
	events, has := notificationUseCase.notificationRepo.GetAll(uid)
	if !has {
		logger.Error("no notifications for the user", uid)
		return nil, false
	}
	for idx := range events {
		var err error
		events[idx].MakeUsr, err = notificationUseCase.userRepo.GetByID(events[idx].MakeUid)
		if err != nil {
			logger.Error(err)
		}
		if events[idx].MetaData.Uid != 0 {
			events[idx].MetaData.Usr, err = notificationUseCase.userRepo.GetByID(events[idx].MetaData.Uid)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return events, true
}

func (notificationUseCase *NotificationUseCase) UpdateAll(uid uint) error {
	if err := notificationUseCase.notificationRepo.UpdateAll(uid); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (notificationUseCase *NotificationUseCase) DeleteAll(uid uint) error {
	if err := notificationUseCase.notificationRepo.DeleteAll(uid); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
