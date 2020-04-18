package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

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
	return usr, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.DB.Where("nickname = ?", nickname).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}
	return usr, nil
}

func (userStore *UserStore) GetBoardsByID(uid uint) ([]models.Board, []models.Board, error) {
	var adminsBoards []models.Board
	usr := &models.User{ID: uid}
	err := userStore.DB.Model(usr).Preload("Admins").Related(&adminsBoards, "Admin").Error
	if err != nil {
		logger.Error(err)
		return nil, nil, errors.ErrUserNotFound
	}
	var membersBoards []models.Board
	err = userStore.DB.Model(usr).Preload("Members").Related(&membersBoards, "Member").Error
	if err != nil {
		logger.Error(err)
		return nil, nil, errors.ErrBoardNotFound
	}
	return adminsBoards, membersBoards, nil
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

func (userStore *UserStore) Update(oldPass []byte, newUser *models.User) error {
	return nil
}

func (userStore *UserStore) Delete(id uint) error {
	if err := userStore.DB.Where("id = ?", id).Delete(models.User{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	return nil
}
