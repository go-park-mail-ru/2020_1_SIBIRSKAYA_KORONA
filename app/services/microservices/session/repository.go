package session

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	Create(sess models.Session) error
	Get(sid string) (uint, error)
	Delete(sid string) error
}
