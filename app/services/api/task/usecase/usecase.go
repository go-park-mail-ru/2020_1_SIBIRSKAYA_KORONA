package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type TaskUseCase struct {
	taskRepo task.Repository
}

func CreateUseCase(taskRepo_ task.Repository) task.UseCase {
	return &TaskUseCase{taskRepo: taskRepo_}
}

func (taskUseCase *TaskUseCase) Create(tsk *models.Task) error {
	return taskUseCase.taskRepo.Create(tsk)
}

func (taskUseCase *TaskUseCase) Get(cid uint, tid uint) (*models.Task, error) {
	tsk, err := taskUseCase.taskRepo.Get(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if tsk.Cid != cid {
		return nil, errors.ErrNoPermission
	}
	return tsk, nil
}

func (taskUseCase *TaskUseCase) Update(newTask models.Task) error {
	err := taskUseCase.taskRepo.Update(newTask)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (taskUseCase *TaskUseCase) Delete(tid uint) error {
	err := taskUseCase.taskRepo.Delete(tid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
