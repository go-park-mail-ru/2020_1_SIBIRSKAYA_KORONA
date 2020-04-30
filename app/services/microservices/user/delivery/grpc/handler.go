package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/golang/protobuf/ptypes/empty"
)

type UserHandler struct {
	UseCase user.UseCase
}

func CreateHandler(useCase user.UseCase) proto.UserServer {
	return &UserHandler{
		UseCase: useCase,
	}
}

func (userHandler *UserHandler) Create(ctx context.Context, mess *proto.UserMess) (*proto.UidMess, error) {
	usr := models.CreateUserFromProto(*mess)
	if usr == nil {
		logger.Error(errors.ErrInternal)
		return &proto.UidMess{Uid: 0}, errors.ErrInternal
	}
	err := userHandler.UseCase.Create(usr)
	if err != nil {
		logger.Error(err)
		return &proto.UidMess{}, err
	}
	return &proto.UidMess{Uid: uint64(usr.ID)}, nil
}

func (userHandler *UserHandler) GetByID(ctx context.Context, mess *proto.UidMess) (*proto.UserMess, error) {
	usr, err := userHandler.UseCase.GetByID(uint(mess.Uid))
	if err != nil {
		logger.Error(err)
		return &proto.UserMess{}, err
	}
	return usr.ToProto(), nil
}

func (userHandler *UserHandler) GetByNickname(ctx context.Context, mess *proto.NicknameMess) (*proto.UserMess, error) {
	usr, err := userHandler.UseCase.GetByNickname(mess.Nickname)
	if err != nil {
		logger.Error(err)
		return &proto.UserMess{}, err
	}
	return usr.ToProto(), nil
}

func (userHandler *UserHandler) CheckPassword(ctx context.Context, mess *proto.CheckPassMess) (*proto.BoolMess, error) {
	ok := userHandler.UseCase.CheckPassword(uint(mess.Uid), mess.Pass)
	return &proto.BoolMess{Ok: ok}, nil
}

func (userHandler *UserHandler) Update(ctx context.Context, mess *proto.UpdateMess) (*empty.Empty, error) {
	usr := models.CreateUserFromProto(*mess.Usr)
	if usr == nil {
		logger.Error(errors.ErrInternal)
		return &empty.Empty{}, errors.ErrInternal
	}
	err := userHandler.UseCase.Update(mess.OldPass, *usr)
	if err != nil {
		logger.Error(err)
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func (userHandler *UserHandler) Delete(ctx context.Context, mess *proto.UidMess) (*empty.Empty, error) {
	err := userHandler.UseCase.Delete(uint(mess.Uid))
	if err != nil {
		logger.Error(err)
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}
