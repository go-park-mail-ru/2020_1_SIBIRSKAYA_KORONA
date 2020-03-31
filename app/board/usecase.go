package board

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"
)

type UseCase interface {
	Create(sid string, board *models.Board) *cstmerr.CustomUsecaseError
	Update(sid string, board *models.Board) *cstmerr.CustomUsecaseError
	Delete(sid string, bid uint) *cstmerr.CustomUsecaseError

	Get(sid string, bid uint) (*models.Board, *cstmerr.CustomUsecaseError)
	GetAll(sid string) ([]models.Board, []models.Board, *cstmerr.CustomUsecaseError)
}
