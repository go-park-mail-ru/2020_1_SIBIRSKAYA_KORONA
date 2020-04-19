package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	pass "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/password"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) user.Repository {
	return &UserStore{DB: db}
}

func (userStore *UserStore) Create(usr *models.User) error {
	if err := userStore.DB.Create(usr).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (userStore *UserStore) GetByID(id uint) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.DB.Where("id = ?", id).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}
	usr.Password = nil
	return usr, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.DB.Where("nickname = ?", nickname).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}
	usr.Password = nil
	return usr, nil
}

func (userStore *UserStore) CheckPassword(uid uint, password []byte) bool {
	var usr models.User
	if err := userStore.DB.Where("id = ?", uid).First(&usr).Error; err != nil {
		logger.Error(err)
		return false
	}
	tmp := pass.CheckPassword(password, usr.Password)
	return tmp
}

func (userStore *UserStore) Update(oldPass []byte, newUser models.User) error {
	var oldUser models.User
	if err := userStore.DB.Where("id = ?", newUser.ID).First(&oldUser).Error; err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	if len(newUser.Password) != 0 {
		if !pass.CheckPassword(oldPass, oldUser.Password) {
			logger.Error(errors.ErrWrongPassword)
			return errors.ErrWrongPassword
		}
		oldUser.Password = pass.HashPasswordGenSalt(newUser.Password)
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
	if err := userStore.DB.Save(oldUser).Error; err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (userStore *UserStore) Delete(id uint) error {
	if err := userStore.DB.Where("id = ?", id).Delete(models.User{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	return nil
}

func (userStore *UserStore) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	var users []models.User
	err := userStore.DB.Limit(limit).Where("nickname LIKE ?", nicknamePart+"%").Find(&users).Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}
	return users, nil
}
