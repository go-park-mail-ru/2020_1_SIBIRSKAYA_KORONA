package http_test

import (
	"net/http"
	test "net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	commentHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/delivery/http"
	commentMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func TestCreateHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentUsecaseMock := commentMocks.NewMockUseCase(ctrl)
	router := echo.New()
	mw := middleware.CreateMiddleware(nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil)
	commentHandler.CreateHandler(router, commentUsecaseMock, mw)
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	htfUseCaseMock := commentMocks.NewMockUseCase(ctrl)
	handler := commentHandler.CommentHandler{UseCase: htfUseCaseMock}

	var testComments models.Comments
	err := faker.FakeData(&testComments)
	assert.NoError(t, err)

	// good
	{
		var uid uint = 10
		var tid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().GetComments(tid, uid).Return(testComments, nil)
		ctx.Set("uid", uid)
		ctx.Set("tid", tid)
		err = handler.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	// error
	{
		var uid uint = 10
		var tid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().GetComments(tid, uid).Return(testComments, errors.ErrInternal)
		ctx.Set("uid", uid)
		ctx.Set("tid", tid)
		err = handler.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusInternalServerError)
	}

}

func TestDelelte(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	htfUseCaseMock := commentMocks.NewMockUseCase(ctrl)
	handler := commentHandler.CommentHandler{UseCase: htfUseCaseMock}

	// good
	{
		var cid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().Delete(cid).Return(nil)
		ctx.Set("comid", cid)
		err := handler.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	// error
	{
		var cid uint = 10
		ctx := echo.New().NewContext(nil, test.NewRecorder())
		htfUseCaseMock.EXPECT().Delete(cid).Return(errors.ErrDbBadOperation)
		ctx.Set("comid", cid)
		err := handler.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusInternalServerError)
	}

}
