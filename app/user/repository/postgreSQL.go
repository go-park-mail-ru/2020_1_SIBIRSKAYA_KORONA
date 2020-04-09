package repository

import (
	"bytes"
	"golang.org/x/crypto/argon2"
	"io"
	"os"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"crypto/rand"
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

func HashPassword(salt, password []byte) []byte {
	hashPass := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	return append(salt, hashPass...)
}

func HashPasswordGenSalt(password []byte) []byte {
	salt := make([]byte, 8)
	rand.Read(salt)
	return HashPassword(salt, password)
}

func CheckPassword(pass, realHashPass []byte) bool {
	salt := realHashPass[0:8]
	return bytes.Equal(HashPassword(salt, pass), realHashPass)
}

func (userStore *UserStore) CheckPasswordByID(uid uint, pass []byte) bool {
	usr := new(models.User)
	if err := userStore.DB.First(&usr, uid).Error; err != nil {
		logger.Error(err)
		return false
	}
	return CheckPassword(pass, usr.Password)
}

func (userStore *UserStore) Create(usr *models.User) error {
	usr.Password = HashPasswordGenSalt(usr.Password)
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

func (userStore *UserStore) Update(oldPass []byte, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error {
	if newUser == nil {
		return errors.ErrInternal
	}
	oldUser := new(models.User)
	if err := userStore.DB.Where("id = ?", newUser.ID).First(&oldUser).Error; err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}
	if oldPass != nil && newUser.Password != nil {
		if !CheckPassword(oldPass, oldUser.Password) {
			logger.Error(errors.ErrWrongPassword)
			return errors.ErrWrongPassword
		}
		oldUser.Password = HashPasswordGenSalt(newUser.Password)
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
			logger.Error(err)
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
	if err := userStore.DB.Where("id = ?", id).Delete(models.User{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
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
