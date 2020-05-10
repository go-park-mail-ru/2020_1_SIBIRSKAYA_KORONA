package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/jinzhu/gorm"
)

type NotificationStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) notification.Repository {
	return &NotificationStore{DB: db}
}

func (notificationStore *NotificationStore) Create(event *models.Event) error {
	err := notificationStore.DB.Create(event).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (notificationStore *NotificationStore) GetAll(uid uint) (models.Events, bool) {
	var events models.Events
	err := notificationStore.DB.Where("uid = ?", uid).Find(&events).Error
	if err != nil {
		logger.Error(err)
		return nil, false
	}
	return events, true
}

func (notificationStore *NotificationStore) UpdateAll(uid uint) error {
	err := notificationStore.DB.Model(models.Events{}).
		Where("uid = ? ", uid).UpdateColumn("is_read", true).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (notificationStore *NotificationStore) DeleteAll(uid uint) error {
	err := notificationStore.DB.Where("uid = ? ", uid).Delete(models.Events{}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}
