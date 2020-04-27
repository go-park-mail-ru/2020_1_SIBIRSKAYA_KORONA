package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type AttachStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) attach.Repository {
	return &AttachStore{DB: db}
}

func (attachStore *AttachStore) Create(attachModel *models.AttachedFile) error {
	err := attachStore.DB.Create(attachModel).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (attachStore *AttachStore) Get(tid uint) (models.AttachedFiles, error) {
	var attachedFiles []models.AttachedFile
	err := attachStore.DB.Model(&models.Task{ID: tid}).Order("id").Related(&attachedFiles, "tid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}

	return attachedFiles, nil
}

func (attachStore *AttachStore) GetByID(fid uint) (*models.AttachedFile, error) {
	attach := new(models.AttachedFile)
	if err := attachStore.DB.First(attach, fid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrFileNotFound
	}

	return attach, nil
}

func (attachStore *AttachStore) Delete(fid uint) error {
	if err := attachStore.DB.Where("id = ?", fid).Delete(models.AttachedFile{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrFileNotFound
	}
	return nil
}
