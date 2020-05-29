package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/labstack/gommon/random"

	"mime/multipart"
)

type AttachUseCase struct {
	attachModelRepo attach.Repository
	attachFileRepo  attach.FileRepository
}

func CreateUseCase(attachModelRepo_ attach.Repository, attachFileRepo_ attach.FileRepository) attach.UseCase {
	return &AttachUseCase{
		attachModelRepo: attachModelRepo_,
		attachFileRepo:  attachFileRepo_,
	}
}

// TODO: переделать под хэш-сумма файла, что позволит не хранить одинаковые файлы на облачном сервисе
// сейчас просто рандомная строчка

func (attachUseCase *AttachUseCase) Create(attachModel *models.AttachedFile, attachFile *multipart.FileHeader) error {
	publicURL, err := attachUseCase.attachFileRepo.UploadFile(attachFile, attachModel)
	if err != nil {
		logger.Error(err)
		return err
	}

	attachModel.URL = publicURL
	attachModel.Name = attachFile.Filename
	attachModel.FileKey = random.String(32, random.Alphabetic, random.Numeric)
	err = attachUseCase.attachModelRepo.Create(attachModel)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (attachUseCase *AttachUseCase) GetByID(tid uint, fid uint) (*models.AttachedFile, error) {
	attach, err := attachUseCase.attachModelRepo.GetByID(fid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if attach.Tid != tid {
		return nil, errors.ErrNoPermission
	}
	return attach, nil
}

func (attachUseCase *AttachUseCase) Get(tid uint) (models.AttachedFiles, error) {
	attachedFiles, err := attachUseCase.attachModelRepo.Get(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return attachedFiles, nil
}

func (attachUseCase *AttachUseCase) Delete(fid uint) error {
	attachFile, err := attachUseCase.attachModelRepo.GetByID(fid)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = attachUseCase.attachFileRepo.DeleteFile(attachFile.FileKey)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = attachUseCase.attachModelRepo.Delete(fid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
