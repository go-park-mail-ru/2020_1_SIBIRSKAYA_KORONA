package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type ItemUseCase struct {
	itemRepo item.Repository
}

func CreateUseCase(itemRepo_ item.Repository) item.UseCase {
	return &ItemUseCase{
		itemRepo: itemRepo_,
	}
}

func (itemUseCase *ItemUseCase) Create(item *models.Item) error {
	err := itemUseCase.itemRepo.Create(item)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// func (itemUseCase *ItemUseCase) Get(tid uint) (models.Checklists, error) {
// 	checklists, err := itemUseCase.itemRepo.Get(tid)
// 	if err != nil {
// 		logger.Error(err)
// 		return nil, err
// 	}
// 	return checklists, nil
// }

func (itemUseCase *ItemUseCase) Update(newItem *models.Item) error {
	err := itemUseCase.itemRepo.Update(newItem)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (itemUseCase *ItemUseCase) Delete(itid uint) error {
	return errors.ErrDbBadOperation
}
