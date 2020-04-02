package usecase

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
)

type UserUseCase struct {
	sessionRepo session.Repository
	userRepo    user.Repository
}

func CreateUseCase(sessionRepo_ session.Repository, userRepo_ user.Repository) user.UseCase {
	return &UserUseCase{
		sessionRepo: sessionRepo_,
		userRepo:    userRepo_,
	}
}

func (userUseCase *UserUseCase) Create(user *models.User, sessionExpires time.Time) (string, error) {
	err := userUseCase.userRepo.Create(user)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	ses := &models.Session{
		SID:     "",
		ID:      user.ID,
		Expires: sessionExpires,
	}

	sid, repoErr := userUseCase.sessionRepo.Create(ses)
	if repoErr != nil {
		logger.Error(repoErr)
		return "", repoErr
	}

	return sid, nil
}

func (userUseCase *UserUseCase) GetByUserKey(userKey string) (*models.User, error) {
	var id uint
	_, err := fmt.Sscan(userKey, &id)
	usr := new(models.User)
	if err == nil {
		usr, err = userUseCase.userRepo.GetByID(id)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	} else {
		usr, err = userUseCase.userRepo.GetByNickname(userKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}
	if usr != nil {
		usr.Password = ""
	}
	return usr, nil
}

func (userUseCase *UserUseCase) GetByCookie(sid string) (*models.User, error) {
	id, has := userUseCase.sessionRepo.Get(sid)
	if !has {
		return nil, errors.ErrSessionNotExist
	}
	usr, err := userUseCase.userRepo.GetByID(id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return usr, nil
}

func (userUseCase *UserUseCase) Update(sid string, oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error {
	if newUser == nil {
		return errors.ErrUserBadMarshall
	}

	id, has := userUseCase.sessionRepo.Get(sid)
	if !has {
		return errors.ErrSessionNotExist
	}

	newUser.ID = id

	repoErr := userUseCase.userRepo.Update(oldPass, newUser, avatarFileDescriptor)
	if repoErr != nil {
		logger.Error(repoErr)
		return repoErr
	}

	return nil
}

func (userUseCase *UserUseCase) Delete(sid string) error {
	id, has := userUseCase.sessionRepo.Get(sid)
	if !has {
		return errors.ErrUserNotExist
	}
	err := userUseCase.sessionRepo.Delete(sid)
	if err != nil {
		return err
	}
	return userUseCase.userRepo.Delete(id)
}
