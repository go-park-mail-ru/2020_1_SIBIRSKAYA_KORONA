package usecase_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/mocks"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/usecase"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"
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

func createRepoMocks(controller *gomock.Controller) (*userMocks.MockRepository, *sessionMocks.MockRepository) {
	userRepoMock := userMocks.NewMockRepository(controller)
	sessionRepoMock := sessionMocks.NewMockRepository(controller)
	return userRepoMock, sessionRepoMock
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	ses := &models.Session{ID: testUser.ID, Expires: 5}

	// good
	sid := "test_sid"
	userRepoMock.EXPECT().GetByNickname(testUser.Nickname).Return(&testUser, nil)
	userRepoMock.EXPECT().CheckPassword(testUser.ID, testUser.Password).Return(true)
	sessionRepoMock.EXPECT().Create(ses).Return(sid, nil)
	resSid, err := sUsecase.Create(&testUser, ses.Expires)
	assert.NoError(t, err)
	assert.Equal(t, sid, resSid)

	// error user not found
	testUser.Nickname += "a"
	userRepoMock.EXPECT().GetByNickname(testUser.Nickname).Return(nil, errors.ErrUserNotFound)
	resSid, err = sUsecase.Create(&testUser, ses.Expires)
	assert.Equal(t, resSid, "")
	assert.Equal(t, err, errors.ErrUserNotFound)

	// error user not found
	testUser.Nickname += "a"
	userRepoMock.EXPECT().GetByNickname(testUser.Nickname).Return(nil, errors.ErrUserNotFound)
	resSid, err = sUsecase.Create(&testUser, ses.Expires)
	assert.Equal(t, resSid, "")
	assert.Equal(t, err, errors.ErrUserNotFound)

	// error user not found
	testUser.ID++
	userRepoMock.EXPECT().GetByNickname(testUser.Nickname).Return(&testUser, nil)
	userRepoMock.EXPECT().CheckPassword(testUser.ID, testUser.Password).Return(false)
	resSid, err = sUsecase.Create(&testUser, ses.Expires)
	assert.Equal(t, resSid, "")
	assert.Equal(t, err, errors.ErrWrongPassword)
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	sid := "test_sid"
	var id uint = 1

	// has
	sessionRepoMock.EXPECT().Get(sid).Return(id, true)
	resId, has := sUsecase.Get(sid)
	assert.Equal(t, resId, id)
	assert.Equal(t, has, true)

	// not has
	sid += "a"
	sessionRepoMock.EXPECT().Get(sid).Return(uint(0), false)
	resId, has = sUsecase.Get(sid)
	assert.Equal(t, resId, uint(0))
	assert.Equal(t, has, false)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	sid := "test_sid"

	// good
	sessionRepoMock.EXPECT().Delete(sid).Return(nil)
	err := sUsecase.Delete(sid)
	assert.NoError(t, err)

	// error
	sid += "a"
	sessionRepoMock.EXPECT().Delete(sid).Return(errors.ErrNoCookie)
	err = sUsecase.Delete(sid)
	assert.Equal(t, err, errors.ErrNoCookie)
}
