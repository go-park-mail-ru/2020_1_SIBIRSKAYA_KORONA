package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type UserStore struct {
	users map[string]*models.User
	mu    sync.Mutex // RWMutex в лекции?
}

func CreateUserStore() user.Repository {
	return &UserStore{
		users: make(map[string]*models.User),
		mu:    sync.Mutex{},
	}
}

func (this *UserStore) AddUser(user *models.User) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.users[user.Nickname] = user
}

func (this *UserStore) GetUser(nickname string) (*models.User, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	user, has := this.users[nickname]
	return user, has
}
