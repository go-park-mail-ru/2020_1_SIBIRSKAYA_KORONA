package repository

import (
	"errors"
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

func (userStore *UserStore) GetByNickname(nickname string) *models.User {
	userData := new(models.User)
	if userStore.DB.Where("nickname = ?", nickname).First(&userData).Error != nil {
		return nil
	}
	userData.Password = ""
	return userData
}

func (userStore *UserStore) Update(oldPass string, newUser *models.User) error {
	oldUser := new(models.User)
	if userStore.DB.First(&oldUser, newUser.ID).Error != nil {
		return nil
	}
	if oldPass != "" && newUser.Password != "" {
		if oldUser.Password != oldPass {
			return errors.New("wrong password")
		}
		oldUser.Password = newUser.Password
	}
	if newUser.Name != "" {
		oldUser.Name = newUser.Name
	}
	if newUser.Surname != "" {
		oldUser.Surname = newUser.Surname
	}
	if newUser.Nickname != "" {
		oldUser.Nickname = newUser.Nickname
	}
	if newUser.Email != "" {
		oldUser.Email = newUser.Email
	}
	if newUser.Img != "" {
		oldUser.Img = newUser.Img
	}
	return userStore.DB.Save(oldUser).Error
}

func (userStore *UserStore) Delete(id uint) error {
	return userStore.DB.Delete(models.User{}, " = ?", id).Error
}
