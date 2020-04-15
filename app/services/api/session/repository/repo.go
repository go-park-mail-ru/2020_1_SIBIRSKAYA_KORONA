package repository

import (
	"context"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/google/uuid"
)

type SessionStore struct {
	clt proto.SessionClient
	ctx context.Context
}

func CreateRepository(clt proto.SessionClient) session.Repository {
	return &SessionStore{
		clt: clt,
		ctx: context.Background(),
	}
}

func (sessionStore *SessionStore) Create(session *models.Session) (string, error) {
	if session == nil {
		return "", errors.ErrInternal
	}
	// TODO: вынести в сервис авторизации генерацию uuid
	session.SID = uuid.New().String()
	params := &proto.CreateMess{
		Sid:        session.SID,
		Uid:        uint32(session.ID),
		Expiration: int32(session.Expires.Unix()),
	}
	res, err := sessionStore.clt.Create(sessionStore.ctx, params)
	// TODO: ошибки
	if err != nil {
		logger.Error(err)
		return "", err
	}
	if res.Error != "" {
		logger.Error(res.Error)
		return "", errors.ErrDbBadOperation
	}
	return session.SID, nil
}

func (sessionStore *SessionStore) Get(sid string) (uint, bool) {
	return 0, sid != ""
}

func (sessionStore *SessionStore) Delete(sid string) error {
	return nil
}
