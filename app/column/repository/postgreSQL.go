package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
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

func (columnStore *ColumnStore) GetTasksByID(cid uint) ([]models.Task, error) {
	var tsks []models.Task
	err := columnStore.DB.Model(&models.Column{ID: cid}).Related(&tsks, "cid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrColNotFound
	}
	return tsks, nil
}
