package session

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type Repository interface {
	Create(session *models.Session) (string, error)
	Get(sid string) (uint, bool)
	Delete(sid string) error
}
