package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ColumnStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) column.Repository {
	return &ColumnStore{DB: db}
}

func (columnStore *ColumnStore) Create(column *models.Column) error {
	err := columnStore.DB.Create(column).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (columnStore *ColumnStore) Get(cid uint) (*models.Column, error) {
	col := new(models.Column)
	if err := columnStore.DB.First(col, cid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrColNotFound
	}
	return col, nil
}

func (columnStore *ColumnStore) Update(newCol models.Column) error {
	oldCol, err := columnStore.Get(newCol.ID)
	if err != nil {
		logger.Error(err)
		return err
	}
	if newCol.Name != "" {
		oldCol.Name = newCol.Name
	}
	if err := columnStore.DB.Save(oldCol).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (columnStore *ColumnStore) GetTasksByID(cid uint) (models.Tasks, error) {
	var tsks []models.Task
	err := columnStore.DB.Model(&models.Column{ID: cid}).Related(&tsks, "cid").Error
	// err := columnStore.DB.Model(&models.Column{ID: cid}).Preload("Members").Related(&tsks, "cid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrColNotFound
	}
	// TODO: попробовать через preload
	// наполняем таску назначенными пользователями и лейблами
	for id := range tsks {
		err := columnStore.DB.Model(tsks[id]).Related(&tsks[id].Members, "Members").Error
		if err != nil {
			logger.Error(err)
			return nil, errors.ErrDbBadOperation
		}
		err = columnStore.DB.Model(tsks[id]).Related(&tsks[id].Labels, "Labels").Error
		if err != nil {
			logger.Error(err)
			return nil, errors.ErrDbBadOperation
		}
	}

	return tsks, nil
}

func (columnStore *ColumnStore) Delete(cid uint) error {
	if err := columnStore.DB.Where("id = ?", cid).Delete(models.Column{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrTaskNotFound
	}
	return nil
}
