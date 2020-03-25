package repository

import (
	"io"
	"os"

	"log"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

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

func (userStore *UserStore) Update(oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) *cstmerr.CustomRepositoryError {
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

	if avatarFileDescriptor != nil {
		urlToSave, err := uploadAvatarToStaticStorage(avatarFileDescriptor, oldUser.Nickname)
		if err != nil {
			log.Println(err)
		} else {
			oldUser.Avatar = urlToSave
		}
	}

	if userStore.DB.Save(oldUser).Error != nil {
		return &cstmerr.CustomRepositoryError{Err: models.ErrDbBadOperation}
	}
	return &cstmerr.CustomRepositoryError{Err: nil}
}

func (userStore *UserStore) Delete(id uint) error {
	err := userStore.DB.Delete(models.User{}, " = ?", id).Error
	return &cstmerr.CustomRepositoryError{Err: err}
}

// не тащим наружу, костыль костылём
func uploadAvatarToStaticStorage(avatarFileDescriptor *multipart.FileHeader, nickname string) (string, error) {
	avatarFile, err := avatarFileDescriptor.Open()
	if err != nil {
		log.Println("Bad avatar file open: ", err)
		return "", err
	}
	defer avatarFile.Close()
	_, format, err := image.DecodeConfig(avatarFile)
	if err != nil {
		log.Println("Bad avatar decoding: ", err)
		return "", err
	}

	avatarFileName := fmt.Sprintf("%s.%s", nickname, format)
	avatarPath := fmt.Sprintf("%s/%s", viper.GetString("frontend.public_dir")+viper.GetString("frontend.avatar_dir"), avatarFileName)
	avatarDst, err := os.Create(avatarPath)
	if err != nil {
		return "", err
	}
	defer avatarDst.Close()

	avatarFile.Seek(0, io.SeekStart)
	_, err = io.Copy(avatarDst, avatarFile)
	if err != nil {
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
