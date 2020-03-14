package session

type Repository interface {
	Create(id uint) (string, error)
	Has(sid string) bool
	Delete(sid string) error
}