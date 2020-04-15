package grpc

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session"
)

type SessionHandler struct {
	useCase session.UseCase
}

func CreateHandler(useCase_ session.UseCase) proto.SessionServer {
	return &SessionHandler{
		useCase: useCase_,
	}
}

func (sessionHandler *SessionHandler) Create(ctx context.Context, mess *proto.CreateMess) (*proto.ErrorMess, error) {
	fmt.Println(mess)
	return nil, nil
}

/*func NewAuthHandler(usecase auth.AuthUsecase) *AuthHandler {
	h := AuthHandler{usecase: usecase}
	h.server = grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
	}))
	auth.RegisterAuthServer(h.server, &h)
	return &h
}*/

/*func (h *AuthHandler) Serve(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Println("listening " + address)
	return h.server.Serve(listener)
}*/