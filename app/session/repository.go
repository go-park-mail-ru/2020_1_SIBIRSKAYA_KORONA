package session

type Repository interface {
	AddSession(nickname string) string
	GetSession(SID string) (string, bool)
	DeleteSession(SID string) error
}