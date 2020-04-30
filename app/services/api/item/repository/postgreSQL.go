package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ItemStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) item.Repository {
	return &ItemStore{DB: db}
}

func (itemStore *ItemStore) Create(itm *models.Item) error {
	err := itemStore.DB.Create(itm).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (itemStore *ItemStore) Update(newItem *models.Item) error {
	var oldItem models.Item
	if err := itemStore.DB.Where("id = ?", newItem.ID).First(&oldItem).Error; err != nil {
		logger.Error(err)
		return errors.ErrItemNotFound
	}
	if newItem.Text != "" {
		oldItem.Text = newItem.Text
	}
	if newItem.IsDone != oldItem.IsDone {
		oldItem.IsDone = newItem.IsDone
	}
	if err := itemStore.DB.Save(oldItem).Error; err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (itemStore *ItemStore) GetByID(itid uint) (*models.Item, error) {
	itm := new(models.Item)
	if err := itemStore.DB.First(itm, itid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrItemNotFound
	}
	return itm, nil
}

func (itemStore *ItemStore) Delete(itid uint) error {
	err := itemStore.DB.Delete(&models.Item{ID: itid}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}
