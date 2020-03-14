package session

type Usecase interface {
	AddSession(nickname string) string
	GetSession(SID string) (string, error)
	DeleteSession(SID string) error
}
