package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type LabelStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) label.Repository {
	return &LabelStore{DB: db}
}

func (labelStore *LabelStore) Create(lbl *models.Label) error {
	err := labelStore.DB.Create(lbl).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (labelStore *LabelStore) Get(lid uint) (*models.Label, error) {
	lbl := new(models.Label)
	if err := labelStore.DB.Where("id = ?", lid).First(lbl).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrLabelNotFound
	}
	return lbl, nil
}

func (labelStore *LabelStore) Update(newLabel models.Label) error {
	oldLabel, err := labelStore.Get(newLabel.ID)
	if err != nil {
		logger.Error(err)
		return err
	}
	if oldLabel.Name != "" {
		oldLabel.Name = newLabel.Name
	}
	if oldLabel.Color != "" {
		oldLabel.Color = newLabel.Color
	}
	if err := labelStore.DB.Save(oldLabel).Error; err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (labelStore *LabelStore) Delete(lid uint) error {
	if err := labelStore.DB.Where("id = ?", lid).Delete(models.Label{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrLabelNotFound
	}
	return nil
}
