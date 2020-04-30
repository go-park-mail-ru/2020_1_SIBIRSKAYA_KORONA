package usecase_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session/mocks"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func createRepoMocks(controller *gomock.Controller) *sessionMocks.MockRepository {
	sessionRepoMock := sessionMocks.NewMockRepository(controller)
	return sessionRepoMock
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock)

	ses := models.Session{ID: 1, SID: "AAA", Expires: 5}

	// good
	sessionRepoMock.EXPECT().Create(ses).Return(nil)
	err := sUsecase.Create(ses)
	assert.NoError(t, err)

	// error
	ses.ID++
	sessionRepoMock.EXPECT().Create(ses).Return(errors.ErrConflict)
	err = sUsecase.Create(ses)
	assert.Equal(t, err, errors.ErrConflict)
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock)

	var id uint = 1
	sid := "test_sid"

	// good
	sessionRepoMock.EXPECT().Get(sid).Return(id, nil)
	resId, err := sUsecase.Get(sid)
	assert.NoError(t, err)
	assert.Equal(t, resId, id)

	// error
	sid += "a"
	sessionRepoMock.EXPECT().Get(sid).Return(uint(0), errors.ErrNoCookie)
	resId, err = sUsecase.Get(sid)
	assert.Equal(t, resId, uint(0))
	assert.Equal(t, errors.ErrNoCookie, err)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock)

	sid := "test_sid"

	// good
	sessionRepoMock.EXPECT().Delete(sid).Return(nil)
	err := sUsecase.Delete(sid)
	assert.NoError(t, err)

	// error
	sid += "a"
	sessionRepoMock.EXPECT().Delete(sid).Return(errors.ErrNoCookie)
	err = sUsecase.Delete(sid)
	assert.Equal(t, errors.ErrNoCookie, err)
}
