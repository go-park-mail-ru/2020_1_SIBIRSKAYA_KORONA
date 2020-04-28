package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/golang/protobuf/ptypes/empty"
)

type SessionHandler struct {
	useCase session.UseCase
}

func CreateHandler(useCase_ session.UseCase) proto.SessionServer {
	return &SessionHandler{
		useCase: useCase_,
	}
}

func (sessionHandler *SessionHandler) Create(ctx context.Context, mess *proto.SessionMess) (*empty.Empty, error) {
	ses := models.CreateSessionFromProto(*mess)
	if ses == nil {
		return &empty.Empty{}, errors.ErrInternal
	}
	err := sessionHandler.useCase.Create(*ses)
	if err != nil {
		logger.Error(err)
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func (sessionHandler *SessionHandler) Get(ctx context.Context, mess *proto.SidMess) (*proto.UidMess, error) {
	uid, err := sessionHandler.useCase.Get(mess.Sid)
	if err != nil {
		logger.Error(err)
		return &proto.UidMess{}, err
	}
	return &proto.UidMess{Uid: uint64(uid)}, nil
}

func (sessionHandler *SessionHandler) Delete(ctx context.Context, mess *proto.SidMess) (*empty.Empty, error) {
	err := sessionHandler.useCase.Delete(mess.Sid)
	if err != nil {
		logger.Error(err)
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}
