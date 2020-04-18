package http

import (
	"context"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	// "github.com/golang/protobuf/ptypes/empty"
)

type UserHandler struct {
	useCase user.UseCase
}

func CreateHandler(useCase user.UseCase) proto.UserServer {
	return &UserHandler{
		useCase: useCase,
	}
}

func (userHandler *UserHandler) Create(ctx context.Context, mess *proto.UserMess) (*proto.UidMess, error) {
	usr := models.CreateUserFromProto(*mess)
	if usr == nil {
		return &proto.UidMess{Uid: 0}, errors.ErrInternal
	}
	err := userHandler.useCase.Create(usr)
	if err != nil {

	}
	return &proto.UidMess{Uid: uint64(usr.ID)}, err
}
