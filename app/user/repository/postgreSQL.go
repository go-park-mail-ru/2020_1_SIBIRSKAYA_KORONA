package repository

import (
	"io"
	"os"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"fmt"
	"mime/multipart"

	// Понятное дело, что заниматься декодированием картинок на бэкенд-сервере плохо,
	// но до появления отдельного решения пока пусть будет так

	"image"
	_ "image/jpeg"
	_ "image/png"
)

type UserStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) user.Repository {
	return &UserStore{DB: db}
}

func (userStore *UserStore) Create(user *models.User) error {
	if err := userStore.DB.Create(user).Error; err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (userStore *UserStore) GetByID(id uint) (*models.User, error) {
	userData := new(models.User)
	if err := userStore.DB.First(&userData, id).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	return userData, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	userData := new(models.User)
	if err := userStore.DB.Where("nickname = ?", nickname).First(&userData).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	return userData, nil
}

func (userStore *UserStore) GetBoardsByID(uid uint) ([]models.Board, []models.Board, error) {
	var adminsBoards, membersBoards []models.Board
	usr := &models.User{ID: uid}
	err := userStore.DB.Model(usr).Preload("Admins").Related(&adminsBoards, "Admin").
		Preload("Members").Related(&membersBoards, "Member").Error
	if err != nil {
		logger.Error(err)
		return nil, nil, errors.ErrDbBadOperation
	}
	return adminsBoards, membersBoards, nil
}

func (userStore *UserStore) Update(oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error {
	if newUser == nil {
		return errors.ErrInternal
	}

	oldUser := new(models.User)

	if err := userStore.DB.First(&oldUser, newUser.ID).Error; err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}

	if oldPass != "" && newUser.Password != "" {
		if oldUser.Password != oldPass {
			return errors.ErrWrongPassword
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

	if avatarFileDescriptor != nil {
		urlToSave, err := uploadAvatarToStaticStorage(avatarFileDescriptor, oldUser.Nickname)
		if err != nil {
			return errors.ErrBadAvatarUpload
		} else {
			oldUser.Avatar = urlToSave
		}
	}

	if err := userStore.DB.Save(oldUser).Error; err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}

	return nil
}

func (userStore *UserStore) Delete(id uint) error {
	if err := userStore.DB.Delete(models.User{}, " = ?", id).Error; err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}

	return nil
}

// не тащим наружу, костыль костылём
func uploadAvatarToStaticStorage(avatarFileDescriptor *multipart.FileHeader, nickname string) (string, error) {
	avatarFile, err := avatarFileDescriptor.Open()
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer avatarFile.Close()
	_, format, err := image.DecodeConfig(avatarFile)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	avatarFileName := fmt.Sprintf("%s.%s", nickname, format)
	avatarPath := fmt.Sprintf("%s/%s", viper.GetString("frontend.public_dir")+viper.GetString("frontend.avatar_dir"), avatarFileName)
	avatarDst, err := os.Create(avatarPath)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer avatarDst.Close()

	avatarFile.Seek(0, io.SeekStart)
	_, err = io.Copy(avatarDst, avatarFile)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	frontStorage := fmt.Sprintf("%s://%s:%s%s",
		viper.GetString("frontend.protocol"),
		viper.GetString("frontend.ip"),
		viper.GetString("frontend.port"),
		viper.GetString("frontend.avatar_dir"))

	avatarStaticUrl := fmt.Sprintf("%s/%s", frontStorage, avatarFileName)
	return avatarStaticUrl, nil
}
