package grpc_test

import (
	"context"
	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/delivery/grpc"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/golang/mock/gomock"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}
	ctx := context.Background()

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	// good
	userUseCaseMock.EXPECT().Create(&testUser).Return(nil)
	res, err := handler.Create(ctx, testUser.ToProto())
	assert.NoError(t, err)
	assert.Equal(t, uint(res.Uid), testUser.ID)

	// error
	testUser.ID++
	userUseCaseMock.EXPECT().Create(&testUser).Return(errors.ErrConflict)
	_, err = handler.Create(ctx, testUser.ToProto())
	assert.Equal(t, err, errors.ErrConflict)
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}
	ctx := context.Background()

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	// good
	userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	res, err := handler.GetByID(ctx, &proto.UidMess{Uid: uint64(testUser.ID)})
	assert.NoError(t, err)
	assert.Equal(t, &testUser, models.CreateUserFromProto(*res))

	// error
	testUser.ID++
	userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
	_, err = handler.GetByID(ctx, &proto.UidMess{Uid: uint64(testUser.ID)})
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestGetByNickname(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}
	ctx := context.Background()

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	// good
	userUseCaseMock.EXPECT().GetByNickname(testUser.Nickname).Return(&testUser, nil)
	res, err := handler.GetByNickname(ctx, &proto.NicknameMess{Nickname: testUser.Nickname})
	assert.NoError(t, err)
	assert.Equal(t, &testUser, models.CreateUserFromProto(*res))

	// error
	testUser.Nickname += "aaa"
	userUseCaseMock.EXPECT().GetByNickname(testUser.Nickname).Return(nil, errors.ErrUserNotFound)
	_, err = handler.GetByNickname(ctx, &proto.NicknameMess{Nickname: testUser.Nickname})
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}
	ctx := context.Background()

	checkPass := proto.CheckPassMess{Uid: 1, Pass: []byte("lovelove")}

	// ok
	userUseCaseMock.EXPECT().CheckPassword(uint(checkPass.Uid), checkPass.Pass).Return(true)
	res, err := handler.CheckPassword(ctx, &checkPass)
	assert.NoError(t, err)
	assert.Equal(t, res.Ok, true)

	// не ok
	checkPass.Uid++
	userUseCaseMock.EXPECT().CheckPassword(uint(checkPass.Uid), checkPass.Pass).Return(false)
	res, err = handler.CheckPassword(ctx, &checkPass)
	assert.NoError(t, err)
	assert.Equal(t, res.Ok, false)
}

// Update(ctx context.Context, mess *proto.UpdateMess) (*empty.Empty, error)
func TestUpdate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}
	ctx := context.Background()

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	updateMess := proto.UpdateMess{OldPass: []byte("lovelove"), Usr: testUser.ToProto()}

	// good
	userUseCaseMock.EXPECT().Update(updateMess.OldPass, testUser).Return(nil)
	_, err = handler.Update(ctx, &updateMess)
	assert.NoError(t, err)

	// error
	testUser.ID++
	updateMess.Usr.Id++
	userUseCaseMock.EXPECT().Update(updateMess.OldPass, testUser).Return(errors.ErrUserNotFound)
	_, err = handler.Update(ctx, &updateMess)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}
	ctx := context.Background()

	var id uint = 1

	// good
	userUseCaseMock.EXPECT().Delete(id).Return(nil)
	_, err := handler.Delete(ctx, &proto.UidMess{Uid: uint64(id)})
	assert.NoError(t, err)

	// error
	id++
	userUseCaseMock.EXPECT().Delete(id).Return(errors.ErrUserNotFound)
	_, err = handler.Delete(ctx, &proto.UidMess{Uid: uint64(id)})
	assert.Equal(t, err, errors.ErrUserNotFound)
}
