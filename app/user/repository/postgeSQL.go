package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) user.Repository {
	return &UserStore{DB: db}
}

func (userStore *UserStore) Create(user *models.User) error {
	return userStore.DB.Create(user).Error
}

func (userStore *UserStore) GetByID(id uint) *models.User {
	userData := new(models.User)
	if userStore.DB.First(&userData, id).Error != nil {
		return nil
	}
	userData.Password = ""
	return userData
}

func (userStore *UserStore) GetByNickName(nickName string) *models.User {
	userData := new(models.User)
	if userStore.DB.Where("nick_name = ?", nickName).First(&userData).Error != nil {
		return nil
	}
	userData.Password = ""
	return userData
}

func (userStore *UserStore) GetAll(id uint) *models.User {
	userData := new(models.User)
	if userStore.DB.First(&userData, id).Error != nil {
		return nil
	}
	return userData
}

func (userStore *UserStore) Update(newUser *models.User) error {
	return userStore.DB.Save(newUser).Error
}

func (userStore *UserStore) Delete(id uint) error {
	return userStore.DB.Delete(models.User{}, " = ?", id).Error
}
