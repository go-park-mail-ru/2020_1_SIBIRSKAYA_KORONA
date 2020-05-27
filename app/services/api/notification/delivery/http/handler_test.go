package http_test

import (
	"net/http"
	test "net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	drelloMiddleware "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	ntfHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/delivery/http"
	ntfMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/mocks"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func TestCreateHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ntfUseCaseMock := ntfMocks.NewMockUseCase(ctrl)

	router := echo.New()
	mw := drelloMiddleware.CreateMiddleware(nil, nil,nil,nil,nil,
		nil,nil,nil,nil, nil,nil)
	ntfHandler.CreateHandler(router, ntfUseCaseMock, mw)
}

func TestGetAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	htfUseCaseMock := ntfMocks.NewMockUseCase(ctrl)
	handler := ntfHandler.NotificationHandler{UseCase: htfUseCaseMock}

	var testEvents models.Events
	err := faker.FakeData(&testEvents)
	assert.NoError(t, err)

	// good
	{
		var uid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().GetAll(uid).Return(testEvents, true)
		ctx.Set("uid", uid)
		err = handler.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	// error
	{
		var uid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().GetAll(uid).Return(testEvents, false)
		ctx.Set("uid", uid)
		err = handler.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

}

func TestUpdateAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	htfUseCaseMock := ntfMocks.NewMockUseCase(ctrl)
	handler := ntfHandler.NotificationHandler{UseCase: htfUseCaseMock}

	// good
	{
		var uid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().UpdateAll(uid).Return(nil)
		ctx.Set("uid", uid)
		err := handler.UpdateAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	// error
	{
		var uid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().UpdateAll(uid).Return(errors.ErrDbBadOperation)
		ctx.Set("uid", uid)
		err := handler.UpdateAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusInternalServerError)
	}
}

func TestDeleteAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	htfUseCaseMock := ntfMocks.NewMockUseCase(ctrl)
	handler := ntfHandler.NotificationHandler{UseCase: htfUseCaseMock}

	// good
	{
		var uid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().DeleteAll(uid).Return(nil)
		ctx.Set("uid", uid)
		err := handler.DeleteAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	// error
	{
		var uid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().DeleteAll(uid).Return(errors.ErrDbBadOperation)
		ctx.Set("uid", uid)
		err := handler.DeleteAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusInternalServerError)
	}
}
