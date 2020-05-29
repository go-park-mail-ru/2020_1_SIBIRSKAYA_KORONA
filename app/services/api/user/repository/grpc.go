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
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	pass "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/password"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
)

type UserStore struct {
	clt    proto.UserClient
	ctx    context.Context
	Config *config.UserConfigController
}

func CreateRepository(clt proto.UserClient, Config_ *config.UserConfigController) user.Repository {
	return &UserStore{
		clt:    clt,
		ctx:    context.Background(),
		Config: Config_,
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
		return errors.ResolveFromRPC(err)
	}
	usr.ID = uint(res.Uid)
	return nil
}

func (userStore *UserStore) GetByID(id uint) (*models.User, error) {
	res, err := userStore.clt.GetByID(userStore.ctx, &proto.UidMess{Uid: uint64(id)})
	if err != nil {
		logger.Error(err)
		return nil, errors.ResolveFromRPC(err)
	}
	return models.CreateUserFromProto(*res), nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	res, err := userStore.clt.GetByNickname(userStore.ctx, &proto.NicknameMess{Nickname: nickname})
	if err != nil {
		logger.Error(err)
		return nil, errors.ResolveFromRPC(err)
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
		urlToSave, errN := userStore.UploadAvatarToStaticStorage(avatarFileDescriptor, newUser.ID)
		if errN != nil {
			logger.Error(errN)
			return errors.ErrBadAvatarUpload
		}
		newUser.Avatar = urlToSave
	}
	_, err = userStore.clt.Update(userStore.ctx, &proto.UpdateMess{OldPass: oldPass, Usr: newUser.ToProto()})
	if err != nil {
		logger.Error(err)
		return errors.ResolveFromRPC(err)
	}
	return nil
}

func (userStore *UserStore) Delete(id uint) error {
	_, err := userStore.clt.Delete(userStore.ctx, &proto.UidMess{Uid: uint64(id)})
	if err != nil {
		logger.Error(err)
		return errors.ResolveFromRPC(err)
	}
	return nil
}

func (userStore *UserStore) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	var users []models.User
	return users, nil
}

func (userStore *UserStore) UploadAvatarToStaticStorage(avatarFileDescriptor *multipart.FileHeader, id uint) (string, error) {
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
	avatarFileName := fmt.Sprintf("%d_%s.%s", id, random.String(8, random.Alphabetic, random.Numeric), format)
	avatarPath := fmt.Sprintf("%s/%s", userStore.Config.GetFrontendAvatarDir(), avatarFileName)
	avatarDst, err := os.Create(avatarPath)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer avatarDst.Close()
	_, errSeek := avatarFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Error(errSeek)
	}
	_, err = io.Copy(avatarDst, avatarFile)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	frontStorage := userStore.Config.GetFrontendStorageURL()
	avatarStaticUrl := fmt.Sprintf("%s/%s", frontStorage, avatarFileName)
	return avatarStaticUrl, nil
}
