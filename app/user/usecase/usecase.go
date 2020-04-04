package usecase

import (
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

func (userUseCase *UserUseCase) GetByID(userID uint) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByID(userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (userUseCase *UserUseCase) GetByNickname(nickname string) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByNickname(nickname)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (userUseCase *UserUseCase) GetBoardsByID(uid uint) ([]models.Board, []models.Board, error) {
	adminsBoard, membersBoard, repoErr := userUseCase.userRepo.GetBoardsByID(uid)
	if repoErr != nil {
		logger.Error(repoErr)
		return nil, nil, repoErr
	}
	return adminsBoard, membersBoard, nil
}

func (userUseCase *UserUseCase) Update(oldPass string, newUser *models.User, avatarFileDescriptor *multipart.FileHeader) error {
	if newUser == nil {
		return errors.ErrUserBadMarshall
	}
	repoErr := userUseCase.userRepo.Update(oldPass, newUser, avatarFileDescriptor)
	if repoErr != nil {
		logger.Error(repoErr)
		return repoErr
	}
	return nil
}

func (userUseCase *UserUseCase) Delete(uid uint, sid string) error {
	err := userUseCase.sessionRepo.Delete(sid)
	if err != nil {
		return err
	}
	return userUseCase.userRepo.Delete(uid)
}
