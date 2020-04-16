package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session"

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

func (sessionHandler *SessionHandler) Create(ctx context.Context, mess *proto.CreateReq) (*empty.Empty, error) {
	err := sessionHandler.useCase.Create(mess.Sid, mess.Uid, mess.Expiration)
	return &empty.Empty{}, err
}

func (sessionHandler *SessionHandler) Get(ctx context.Context, mess *proto.GetReq) (*proto.GetResp, error) {
	uid, err := sessionHandler.useCase.Get(mess.Sid)
	return &proto.GetResp{Uid: uint32(uid)}, err
}

func (sessionHandler *SessionHandler) Delete(ctx context.Context, mess *proto.DeleteReq) (*empty.Empty, error) {
	err := sessionHandler.useCase.Delete(mess.Sid)
	return &empty.Empty{}, err
}
