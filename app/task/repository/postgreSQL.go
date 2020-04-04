package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TaskStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) task.Repository {
	return &TaskStore{DB: db}
}

func (taskStore *TaskStore) Create(tsk *models.Task) error {
	err := taskStore.DB.Create(tsk).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (taskStore *TaskStore) Get(tid uint) (*models.Task, error) {
	tsk := new(models.Task)
	if err := taskStore.DB.First(tsk, tid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	return tsk, nil
}
