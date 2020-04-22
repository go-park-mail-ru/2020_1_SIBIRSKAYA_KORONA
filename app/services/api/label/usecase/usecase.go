package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type LabelUseCase struct {
	labelRepo label.Repository
}

func CreateUseCase(labelRepo_ label.Repository) label.UseCase {
	return &LabelUseCase{
		labelRepo: labelRepo_,
	}
}

func (labelUseCase *LabelUseCase) Create(label *models.Label) error {
	err := labelUseCase.labelRepo.Create(label)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (labelUseCase *LabelUseCase) Get(bid uint, lid uint) (*models.Label, error) {
	lbl, err := labelUseCase.labelRepo.Get(lid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if lbl.Bid != bid {
		return nil, errors.ErrNoPermission
	}
	return lbl, nil
}

func (labelUseCase *LabelUseCase) Update(lbl models.Label) error {
	err := labelUseCase.labelRepo.Update(lbl)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (labelUseCase *LabelUseCase) Delete(lid uint) error {
	err := labelUseCase.labelRepo.Delete(lid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (labelUseCase *LabelUseCase) AddLabelOnTask(lid uint, tid uint) error {
	err := labelUseCase.labelRepo.AddLabelOnTask(lid, tid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (labelUseCase *LabelUseCase) RemoveLabelFromTask(lid uint, tid uint) error {
	err := labelUseCase.labelRepo.RemoveLabelFromTask(lid, tid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
