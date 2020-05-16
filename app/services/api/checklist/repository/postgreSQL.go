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
	err := checklistStore.DB.Model(&models.Task{ID: tid}).Order("id").Related(&chlists, "tid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	for id := range chlists {
		err := checklistStore.DB.Model(chlists[id]).Order("id").Related(&chlists[id].Items, "clid").Error
		if err != nil {
			logger.Error(err)
			return nil, errors.ErrDbBadOperation
		}
	}
	return chlists, nil
}

func (checklistStore *ChecklistStore) GetByID(clid uint) (*models.Checklist, error) {
	chlist := new(models.Checklist)
	if err := checklistStore.DB.First(chlist, clid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrChecklistNotFound
	}
	err := checklistStore.DB.Model(&chlist).Related(&chlist.Items, "clid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	return chlist, nil
}

func (checklistStore *ChecklistStore) Delete(clid uint) error {
	var items []models.Item
	err := checklistStore.DB.Model(&models.Checklist{ID: clid}).Related(&items, "clid").Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	for id := range items {
		err = checklistStore.DB.Delete(&models.Item{ID: items[id].ID}).Error
		if err != nil {
			logger.Error(err)
			return errors.ErrDbBadOperation
		}
	}
	err = checklistStore.DB.Delete(&models.Checklist{ID: clid}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (checklistStore *ChecklistStore) Update(checklist *models.Checklist) error {
	return errors.ErrDbBadOperation
}
