package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"
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
	return userData
}

func (userStore *UserStore) GetByNickname(nickname string) *models.User {
	userData := new(models.User)
	if userStore.DB.Where("nickname = ?", nickname).First(&userData).Error != nil {
		return nil
	}
	return userData
}

func (userStore *UserStore) Update(oldPass string, newUser *models.User) *cstmerr.CustomRepositoryError {
	if newUser == nil {
		return &cstmerr.CustomRepositoryError{Err: models.ErrInternal}
	}
	oldUser := new(models.User)
	if userStore.DB.First(&oldUser, newUser.ID).Error != nil {
		return &cstmerr.CustomRepositoryError{Err: models.ErrDbBadOperation}
	}
	if oldPass != "" && newUser.Password != "" {
		if oldUser.Password != oldPass {
			return &cstmerr.CustomRepositoryError{Err: models.ErrWrongPassword}
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
	if newUser.Avatar != "" {
		oldUser.Avatar = newUser.Avatar
	}
	// Обернуть ошибки из базы, чтобы можно было логировать
	if userStore.DB.Save(oldUser).Error != nil {
		return &cstmerr.CustomRepositoryError{Err: models.ErrDbBadOperation}
	}
	return &cstmerr.CustomRepositoryError{Err: nil}
}

func (userStore *UserStore) Delete(id uint) error {
	err := userStore.DB.Delete(models.User{}, " = ?", id).Error
	return &cstmerr.CustomRepositoryError{Err: err}
}
