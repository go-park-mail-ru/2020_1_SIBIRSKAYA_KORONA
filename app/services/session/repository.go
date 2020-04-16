package session

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/session_repo_mock.go
type Repository interface {
	Create(sid string, uid uint32, expiration int32) error
	Get(sid string) (uint, error)
	Delete(sid string) error
}
