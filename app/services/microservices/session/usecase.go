package session

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/session_usecase_mock.go
type UseCase interface {
	Create(sess models.Session) error
	Get(sid string) (uint, error)
	Delete(sid string) error
}
