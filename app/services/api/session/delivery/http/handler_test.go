package http_test

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"net/http"
	test "net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/delivery/http"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/mocks"
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

func GetContexFromJSON(method, path string) echo.Context {
	request := test.NewRequest(method, path, nil)
	return echo.New().NewContext(request, test.NewRecorder())
}

func TestLogIn(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCaseMock := sessionMocks.NewMockUseCase(ctrl)
	handler := sessionHandler.SessionHandler{UseCase: sessionUseCaseMock}

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	// good
	{
		body, err := testUser.MarshalJSON()
		assert.NoError(t, err)
		ctx := GetContexFromJSON(echo.POST, "/session")
		ctx.Set("body", body)
		sessionUseCaseMock.EXPECT().Create(&testUser, gomock.Any()).Return("test_sid", nil)

		err = handler.LogIn(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}
	// error
	{
		testUser.ID++
		body, err := testUser.MarshalJSON()
		assert.NoError(t, err)
		ctx := GetContexFromJSON(echo.POST, "/session")
		ctx.Set("body", body)
		sessionUseCaseMock.EXPECT().Create(&testUser, gomock.Any()).Return("", errors.ErrConflict)

		err = handler.LogIn(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusConflict)
	}
}

func TestLogOut(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCaseMock := sessionMocks.NewMockUseCase(ctrl)
	handler := sessionHandler.SessionHandler{UseCase: sessionUseCaseMock}

	sid := "test_sid"
	// good
	{
		ctx := GetContexFromJSON(echo.DELETE, "/session")
		ctx.Set("sid", sid)
		sessionUseCaseMock.EXPECT().Delete(sid).Return(nil)

		err := handler.LogOut(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}
	// error
	{
		ctx := GetContexFromJSON(echo.DELETE, "/session")
		ctx.Set("sid", sid)
		sessionUseCaseMock.EXPECT().Delete(sid).Return(errors.ErrNoCookie)

		err := handler.LogOut(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusForbidden)
	}
}

func TestToken(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionUseCaseMock := sessionMocks.NewMockUseCase(ctrl)
	handler := sessionHandler.SessionHandler{UseCase: sessionUseCaseMock}

	ctx := GetContexFromJSON(echo.DELETE, "/session")
	sid := "test_sid"
	ctx.Set("sid", sid)

	err := handler.Token(ctx)
	assert.NoError(t, err)
	assert.Equal(t, ctx.Response().Status, http.StatusOK)
}
