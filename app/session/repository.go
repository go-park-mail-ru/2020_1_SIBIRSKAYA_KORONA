package session

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/session_repo_mock.go
type Repository interface {
	Create(session *models.Session) (string, error)
	Get(sid string) (uint, bool)
	Delete(sid string) error
}
