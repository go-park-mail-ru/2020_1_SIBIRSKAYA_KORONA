package usecase_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/mocks"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func createRepoMocks(controller *gomock.Controller) *userMocks.MockRepository {
	userRepoMock := userMocks.NewMockRepository(controller)
	return userRepoMock
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := createRepoMocks(ctrl)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	uUsecase := userUseCase.CreateUseCase(userRepoMock)

	// good
	userRepoMock.EXPECT().Create(&testUser).Return(nil)
	err = uUsecase.Create(&testUser)
	assert.NoError(t, err)

	// error
	testUser.ID++
	userRepoMock.EXPECT().Create(&testUser).Return(errors.ErrConflict)
	err = uUsecase.Create(&testUser)
	assert.Equal(t, err, errors.ErrConflict)
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := createRepoMocks(ctrl)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	uUsecase := userUseCase.CreateUseCase(userRepoMock)

	// good
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	usr, err := uUsecase.GetByID(testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, usr, &testUser)

	// error
	testUser.ID++
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
	_, err = uUsecase.GetByID(testUser.ID)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestGetByNickname(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := createRepoMocks(ctrl)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	uUsecase := userUseCase.CreateUseCase(userRepoMock)

	// good
	userRepoMock.EXPECT().GetByNickname(testUser.Nickname).Return(&testUser, nil)
	usr, err := uUsecase.GetByNickname(testUser.Nickname)
	assert.NoError(t, err)
	assert.Equal(t, usr, &testUser)

	// error
	testUser.Nickname += "aa"
	userRepoMock.EXPECT().GetByNickname(testUser.Nickname).Return(nil, errors.ErrUserNotFound)
	_, err = uUsecase.GetByNickname(testUser.Nickname)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := createRepoMocks(ctrl)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	uUsecase := userUseCase.CreateUseCase(userRepoMock)

	oldPass := []byte("oldPass")
	userRepoMock.EXPECT().Update(oldPass, testUser).Return(nil)
	err = uUsecase.Update(oldPass, testUser)
	assert.NoError(t, err)

	// error not found
	testUser.ID++
	userRepoMock.EXPECT().Update(nil, testUser).Return(errors.ErrUserNotFound)
	err = uUsecase.Update(nil, testUser)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := createRepoMocks(ctrl)

	var id uint = 1
	pass := []byte("12345678")
	uUsecase := userUseCase.CreateUseCase(userRepoMock)

	// ok
	userRepoMock.EXPECT().CheckPassword(id, pass).Return(true)
	ok := uUsecase.CheckPassword(id, pass)
	assert.Equal(t, ok, true)

	// не ok
	id++
	userRepoMock.EXPECT().CheckPassword(id, pass).Return(false)
	ok = uUsecase.CheckPassword(id, pass)
	assert.Equal(t, ok, false)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := createRepoMocks(ctrl)

	var id uint = 1
	uUsecase := userUseCase.CreateUseCase(userRepoMock)

	// good
	userRepoMock.EXPECT().Delete(id).Return(nil)
	err := uUsecase.Delete(id)
	assert.NoError(t, err)

	// err not found
	id++
	userRepoMock.EXPECT().Delete(id).Return(errors.ErrUserNotFound)
	err = uUsecase.Delete(id)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestGetUsersByNicknamePart(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := createRepoMocks(ctrl)
	uUsecase := userUseCase.CreateUseCase(userRepoMock)

	var testUsers models.Users
	err := faker.FakeData(&testUsers)
	assert.NoError(t, err)

	nicknamePart := "aaa"
	var limit uint = 3

	// good
	userRepoMock.EXPECT().GetUsersByNicknamePart(nicknamePart, limit).Return(testUsers, nil)
	usrs, err := uUsecase.GetUsersByNicknamePart(nicknamePart, limit)
	assert.NoError(t, err)
	assert.Equal(t, len(testUsers), len(usrs))
	for idx := range testUsers {
		assert.Equal(t, usrs[idx], testUsers[idx])
	}

	// error not found
	nicknamePart += "a"
	limit--
	userRepoMock.EXPECT().GetUsersByNicknamePart(nicknamePart, limit).Return(nil, errors.ErrUserNotFound)
	_, err = uUsecase.GetUsersByNicknamePart(nicknamePart, limit)
	assert.Equal(t, err, errors.ErrUserNotFound)
}
