package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"

	"github.com/jinzhu/gorm"
)

type NotificationStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) notification.Repository {
	return &NotificationStore{DB: db}
}

func (notificationStore *NotificationStore) Add(evnt models.Event) error {
	return nil
}

func (notificationStore NotificationStore) Pop(uid uint) (models.Events, bool) {
	event := models.Event{Message: "keklol"}
	return models.Events{event}, true
}
