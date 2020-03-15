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
}

func (userStore *UserStore) GetByID(id uint) *models.User {
}

func (userStore *UserStore) GetByNickName(nickName string) *models.User {
}

func (userStore *UserStore)  HasNickName(nickName string) bool {
}

func (userStore *UserStore) GetAll(id uint) *models.User{
}

func (userStore *UserStore) Update(newUser *models.User) error {
}

func (userStore *UserStore) Delete(id uint) error {
}
