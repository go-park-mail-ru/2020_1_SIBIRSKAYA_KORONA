package session

import (
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type UseCase interface {
	Create(user *models.User, sessionExpires time.Time) (string, error)
	Has(sid string) bool
	Delete(sid string) error
}
