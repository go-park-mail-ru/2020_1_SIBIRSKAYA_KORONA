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
		return errors.ErrConflict
	}
	return nil
}

func (taskStore *TaskStore) Get(tid uint) (*models.Task, error) {
	tsk := new(models.Task)
	if err := taskStore.DB.First(tsk, tid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTaskNotFound
	}
	return tsk, nil
}

func (taskStore *TaskStore) Update(newTask models.Task) error {
	oldTask, err := taskStore.Get(newTask.ID)
	if err != nil {
		logger.Error(err)
		return err
	}
	if newTask.Name != "" {
		oldTask.Name = newTask.Name
	}
	if newTask.About != "" {
		oldTask.About = newTask.About
	}
	if newTask.Level != 0 {
		oldTask.Level = newTask.Level
	}
	if newTask.Deadline != "" {
		oldTask.Deadline = newTask.Deadline
	}
	if newTask.Pos != 0 {
		oldTask.Pos = newTask.Pos
	}
	if newTask.Cid != 0 {
		oldTask.Cid = newTask.Cid
	}
	if err := taskStore.DB.Save(oldTask).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (taskStore *TaskStore) Delete(tid uint) error {
	if err := taskStore.DB.Where("id = ?", tid).Delete(models.Task{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrTaskNotFound
	}
	return nil
}
