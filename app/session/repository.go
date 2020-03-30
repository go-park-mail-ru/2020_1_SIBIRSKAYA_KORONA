package session

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"
)

type Repository interface {
	Create(session *models.Session) (string, *cstmerr.RepoError)
	Get(sid string) (uint, bool)
	Delete(sid string) error
}
