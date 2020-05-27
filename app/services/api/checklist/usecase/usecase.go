package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type ChecklistUseCase struct {
	checklistRepo checklist.Repository
	itemRepo      item.Repository
}

func CreateUseCase(checklistRepo_ checklist.Repository, itemRepo_ item.Repository) checklist.UseCase {
	return &ChecklistUseCase{
		checklistRepo: checklistRepo_,
		itemRepo:      itemRepo_,
	}
}

func (checklistUseCase *ChecklistUseCase) Create(chlist *models.Checklist) error {
	err := checklistUseCase.checklistRepo.Create(chlist)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (checklistUseCase *ChecklistUseCase) Get(tid uint) (models.Checklists, error) {
	checklists, err := checklistUseCase.checklistRepo.Get(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return checklists, nil
}

func (checklistUseCase *ChecklistUseCase) GetByID(tid uint, clid uint) (*models.Checklist, error) {
	checklist, err := checklistUseCase.checklistRepo.GetByID(clid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if checklist.Tid != tid {
		return nil, errors.ErrNoPermission
	}
	return checklist, nil
}

func (checklistUseCase *ChecklistUseCase) Update(chlist *models.Checklist) error {
	return errors.ErrDbBadOperation
}

func (checklistUseCase *ChecklistUseCase) Delete(clid uint) error {
	err := checklistUseCase.checklistRepo.Delete(clid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
