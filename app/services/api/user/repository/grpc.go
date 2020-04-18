package repository

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/spf13/viper"
)

type UserStore struct {
	clt proto.UserClient
	ctx context.Context
}

func CreateRepository(clt proto.UserClient) user.Repository {
	return &UserStore{
		clt: clt,
		ctx: context.Background(),
	}
}

func (userStore *UserStore) Create(usr *models.User) error {
	if usr == nil {
		return errors.ErrInternal
	}
	mess := usr.ToProto()
	if mess == nil {
		return errors.ErrInternal
	}
	res, err := userStore.clt.Create(userStore.ctx, mess)
	if err != nil {
		logger.Error(err)
		return err
	}
	usr.ID = uint(res.Uid)
	return nil
}

func (userStore *UserStore) CheckPasswordByID(uid uint, pass []byte) bool {
	return false
}

func (userStore *UserStore) GetByID(id uint) (*models.User, error) {
	usr := new(models.User)
	return usr, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	usr := new(models.User)
	return usr, nil
}

// TODO: вынести в доски
func (userStore *UserStore) GetBoardsByID(uid uint) ([]models.Board, []models.Board, error) {
	var adminsBoards []models.Board
	var membersBoards []models.Board
	return adminsBoards, membersBoards, nil
}

func (userStore *UserStore) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	var users []models.User
	return users, nil
}

func (userStore *UserStore) Update(oldPass []byte, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error {
	return nil
}

func (userStore *UserStore) Delete(id uint) error {
	return nil
}

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
	publicDirPath, exists := os.LookupEnv("DRELLO_PUBLIC_DIR")
	if !exists {
		logger.Error("DRELLO_PUBLIC_DIR environment variable not exist")
	}
	avatarPath := fmt.Sprintf("%s/%s", publicDirPath+viper.GetString("frontend.avatar_dir"), avatarFileName)
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
