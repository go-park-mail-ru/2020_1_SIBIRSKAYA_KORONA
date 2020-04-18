package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"mime/multipart"
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

func (userUseCase *UserUseCase) Create(user *models.User, sessionExpires int64) (string, error) {
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
	sid, err := userUseCase.sessionRepo.Create(ses)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return sid, nil
}

func (userUseCase *UserUseCase) GetByID(uid uint) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByID(uid)
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
	adminsBoard, membersBoard, err := userUseCase.userRepo.GetBoardsByID(uid)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}
	return adminsBoard, membersBoard, nil
}

func (userUseCase *UserUseCase) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	users, err := userUseCase.userRepo.GetUsersByNicknamePart(nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return users, nil
}

func (userUseCase *UserUseCase) Update(oldPass []byte, newUser models.User, avatarFileDescriptor *multipart.FileHeader) error {
	err := userUseCase.userRepo.Update(oldPass, newUser, avatarFileDescriptor)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (userUseCase *UserUseCase) Delete(uid uint, sid string) error {
	err := userUseCase.sessionRepo.Delete(sid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return userUseCase.userRepo.Delete(uid)
}
