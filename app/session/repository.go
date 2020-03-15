package session

type Repository interface {
	Create(id uint) (string, error)
    Get(sid string) (uint, bool)
    Delete(sid string) error
}