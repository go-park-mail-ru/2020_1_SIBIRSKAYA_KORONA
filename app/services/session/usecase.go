package session

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/session_usecase_mock.go

type UseCase interface {
	Create(sid string, uid uint32, expiration int32) error
	Get(sid string) (uint, error)
	Delete(sid string) error
}
