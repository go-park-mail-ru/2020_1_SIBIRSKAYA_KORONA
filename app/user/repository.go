package user

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

type Repository interface {
	AddUser(user *models.User)
	GetUser(nickname string) (*models.User, bool)

	// Влажные мечты
	// GetByName(nickname string) (*models.User, error)
	// GetAllByName(nickname string) (*models.User, error)
	// Create(user *User) error
	// Update(user *User) error
	// Delete(user *User) error
}
