package checklist

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(chlist *models.Checklist) error
	Get(tid uint) (models.Checklists, error)
	Update(chlist *models.Checklist) error
	Delete(clid uint) error
	GetByID(clid uint) (*models.Checklist, error)
}
