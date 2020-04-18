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
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/password"

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
	usr.Password = pass.HashPasswordGenSalt(usr.Password)
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

func (userStore *UserStore) GetByID(id uint) (*models.User, error) {
	res, err := userStore.clt.GetByID(userStore.ctx, &proto.UidMess{Uid: uint64(id)})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return models.CreateUserFromProto(*res), nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	res, err := userStore.clt.GetByNickname(userStore.ctx, &proto.NicknameMess{Nickname: nickname})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return models.CreateUserFromProto(*res), nil
}

func (userStore *UserStore) CheckPassword(uid uint, pass []byte) bool {
	res, err := userStore.clt.CheckPassword(userStore.ctx, &proto.CheckPassMess{Uid: uint64(uid), Pass: pass})
	if err != nil {
		logger.Error(err)
		return false
	}
	return res.Ok
}

func (userStore *UserStore) Update(oldPass []byte, newUser models.User, avatarFileDescriptor *multipart.FileHeader) error {
	_, err := userStore.GetByID(newUser.ID)
	if err != nil {
		logger.Error(err)
		return err
	}
	if avatarFileDescriptor != nil {
		urlToSave, err := UploadAvatarToStaticStorage(avatarFileDescriptor, newUser.ID)
		if err != nil {
			logger.Error(err)
			return errors.ErrBadAvatarUpload
		}
		newUser.Avatar = urlToSave
	}
	_, err = userStore.clt.Update(userStore.ctx, &proto.UpdateMess{OldPass: oldPass, Usr: newUser.ToProto()})
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (userStore *UserStore) Delete(id uint) error {
	_, err := userStore.clt.Delete(userStore.ctx, &proto.UidMess{Uid: uint64(id)})
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (userStore *UserStore) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	var users []models.User
	return users, nil
}

// TODO: вынести в доски
func (userStore *UserStore) GetBoardsByID(uid uint) ([]models.Board, []models.Board, error) {
	var adminsBoards []models.Board
	var membersBoards []models.Board
	return adminsBoards, membersBoards, nil
}

func UploadAvatarToStaticStorage(avatarFileDescriptor *multipart.FileHeader, id uint) (string, error) {
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
	avatarFileName := fmt.Sprintf("%d.%s", id, format)
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
