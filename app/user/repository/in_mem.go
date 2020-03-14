package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
)

type UserStore struct {
	users map[uint]*models.User
	mu    sync.Mutex
}

func CreateRepository() user.Repository {
	return &UserStore{
		users: make(map[uint]*models.User),
		mu:    sync.Mutex{},
	}
}

func (userStore *UserStore) Create(user *models.User) error {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()
	user.ID = uint(len(userStore.users))
	userStore.users[user.ID] = user
	return nil
}

func (userStore *UserStore) Get(id uint) *models.User {
	tmp := userStore.GetAll(id)
	if tmp == nil {
		return tmp
	}
	pubUser := *tmp
	pubUser.Password = ""
	return &pubUser
}

func (userStore *UserStore) GetAll(id uint) *models.User{
	userStore.mu.Lock()
	defer userStore.mu.Unlock()
	u, _ := userStore.users[id]
	return u
}

func (userStore *UserStore) Update(newUser *models.User) error {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()
	userStore.users[newUser.ID] = newUser
	return nil
}

func (userStore *UserStore) Delete(id uint) error {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()
	delete(userStore.users, id)
	return nil
}
