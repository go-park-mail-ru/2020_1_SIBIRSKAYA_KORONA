package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ChecklistStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) checklist.Repository {
	return &ChecklistStore{DB: db}
}

func (checklistStore *ChecklistStore) Create(chlist *models.Checklist) error {
	err := checklistStore.DB.Create(chlist).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (checklistStore *ChecklistStore) Get(tid uint) (models.Checklists, error) {
	var chlists []models.Checklist
	err := checklistStore.DB.Model(&models.Task{ID: tid}).Related(&chlists, "tid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}

	for id := range chlists {
		err := checklistStore.DB.Model(chlists[id]).Related(&chlists[id].Items, "clid").Error
		if err != nil {
			logger.Error(err)
			return nil, errors.ErrDbBadOperation
		}
	}
	return chlists, nil

}

func (checklistStore *ChecklistStore) Delete(clid uint) error {
	err := checklistStore.DB.Delete(&models.Checklist{ID: clid}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrBoardNotFound
	}
	return nil
}

func (checklistStore *ChecklistStore) Update(checklist *models.Checklist) error {
	return errors.ErrDbBadOperation
}
