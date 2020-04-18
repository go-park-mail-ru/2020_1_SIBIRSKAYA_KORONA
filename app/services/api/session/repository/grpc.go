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
	session.SID = uuid.New().String()
	mess := session.ToProto()
	if mess == nil {
		return "", errors.ErrInternal
	}
	_, err := sessionStore.clt.Create(sessionStore.ctx, mess)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return session.SID, nil
}

func (sessionStore *SessionStore) Get(sid string) (uint, bool) {
	res, err := sessionStore.clt.Get(sessionStore.ctx, &proto.SidMess{Sid: sid})
	if err != nil {
		logger.Error(err)
		return 0, false
	}
	return uint(res.Uid), true
}

func (sessionStore *SessionStore) Delete(sid string) error {
	_, err := sessionStore.clt.Delete(sessionStore.ctx, &proto.SidMess{Sid: sid})
	if err != nil {
		logger.Error(err)
	}
	return err
}
