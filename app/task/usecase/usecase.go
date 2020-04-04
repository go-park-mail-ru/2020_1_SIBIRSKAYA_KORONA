package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task"
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
