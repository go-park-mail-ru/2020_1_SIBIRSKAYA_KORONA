package checklist

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type UseCase interface {
	Create(chlist *models.Checklist) error
	Get(tid uint) (models.Checklists, error)
	Update(chlist *models.Checklist) error
	Delete(clid uint) error
	GetByID(tid uint, clid uint) (*models.Checklist, error)
}