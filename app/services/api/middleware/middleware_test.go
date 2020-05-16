package middleware_test

import (
	"os"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

// func TestCORS(t *testing.T) {
// 	 t.Skip()
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sessionUsecaseMock := sessionMocks.NewMockUseCase(ctrl)
// 	boardUsecaseMock := boardMocks.NewMockUseCase(ctrl)
// 	columnUsecaseMock := columnMocks.NewMockUseCase(ctrl)
// 	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)

// 	middle := middleware.CreateMiddleware(sessionUsecaseMock, boardUsecaseMock, columnUsecaseMock, taskUsecaseMock)

// 	e := echo.New()
// 	request := test.NewRequest(echo.GET, "/", nil)
// 	response := test.NewRecorder()
// 	context := e.NewContext(request, response)

// 	executableHandler := middle.CORS(echo.HandlerFunc(func(c echo.Context) error {
// 		return c.NoContent(http.StatusOK)
// 	}))

// 	err := executableHandler(context)
// 	require.NoError(t, err)
// 	assert.Equal(t, "http://localhost:5757", response.Header().Get("Access-Control-Allow-Origin"))
// }

// func TestPanicProcess(t *testing.T) {
// 	// t.Skip()
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sessionUsecaseMock := sessionMocks.NewMockUseCase(ctrl)
// 	boardUsecaseMock := boardMocks.NewMockUseCase(ctrl)
// 	columnUsecaseMock := columnMocks.NewMockUseCase(ctrl)
// 	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)

// 	middle := middleware.CreateMiddleware(sessionUsecaseMock, boardUsecaseMock, columnUsecaseMock, taskUsecaseMock)

// 	e := echo.New()
// 	request := test.NewRequest(echo.GET, "/settings", nil)
// 	response := test.NewRecorder()
// 	context := e.NewContext(request, response)

// 	panicHandler := echo.HandlerFunc(func(c echo.Context) error {
// 		if 2+2 == 4 {
// 			panic("big panic")
// 		}
// 		return c.NoContent(http.StatusOK)
// 	})

// 	processedPanicHandler := middle.ProcessPanic(panicHandler)

// 	err := processedPanicHandler(context)
// 	require.NoError(t, err)
// }

// func TestCheckAuth(t *testing.T) {
// 	// t.Skip()
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sessionUsecaseMock := sessionMocks.NewMockUseCase(ctrl)
// 	boardUsecaseMock := boardMocks.NewMockUseCase(ctrl)
// 	columnUsecaseMock := columnMocks.NewMockUseCase(ctrl)
// 	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)

// 	middle := middleware.CreateMiddleware(sessionUsecaseMock, boardUsecaseMock, columnUsecaseMock, taskUsecaseMock)

// 	e := echo.New()

// 	noCookieRequest := test.NewRequest(echo.GET, "/", nil)
// 	noCookieResponse := test.NewRecorder()
// 	noCookieContext := e.NewContext(noCookieRequest, noCookieResponse)

// 	withCookieRequest := test.NewRequest(echo.GET, "/", nil)
// 	withCookieResponse := test.NewRecorder()
// 	withCookieContext := e.NewContext(withCookieRequest, withCookieResponse)

// 	cookie := http.Cookie{Name: "session_id", Value: "check_cookie"}
// 	withCookieContext.Request().AddCookie(&cookie)

// 	sessionUsecaseMock.EXPECT().
// 		Get("check_cookie").
// 		Return(uint(10), true)

// 	testHandler := echo.HandlerFunc(func(context echo.Context) error {
// 		return context.NoContent(http.StatusOK)
// 	})

// 	executableHandler := middle.CheckAuth(testHandler)

// 	noCookieErr := executableHandler(noCookieContext)
// 	require.NoError(t, noCookieErr)
// 	assert.Equal(t, noCookieContext.Response().Status, http.StatusUnauthorized)

// 	withCookieErr := executableHandler(withCookieContext)
// 	require.NoError(t, withCookieErr)
// 	assert.Equal(t, withCookieContext.Response().Status, http.StatusOK)

// 	assert.Equal(t, withCookieContext.Get("uid").(uint), uint(10))
// 	assert.Equal(t, withCookieContext.Get("sid").(string), "check_cookie")

// }

// func TestCSRFmiddle(t *testing.T) {
// 	// t.Skip()
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sessionUsecaseMock := sessionMocks.NewMockUseCase(ctrl)
// 	boardUsecaseMock := boardMocks.NewMockUseCase(ctrl)
// 	columnUsecaseMock := columnMocks.NewMockUseCase(ctrl)
// 	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)

// 	middle := middleware.CreateMiddleware(sessionUsecaseMock, boardUsecaseMock, columnUsecaseMock, taskUsecaseMock)

// 	e := echo.New()

// 	noTokenRequest := test.NewRequest(echo.GET, "/", nil)
// 	noTokenResponse := test.NewRecorder()
// 	noTokenContext := e.NewContext(noTokenRequest, noTokenResponse)

// 	withTokenRequest := test.NewRequest(echo.GET, "/", nil)
// 	withTokenResponse := test.NewRecorder()
// 	withTokenContext := e.NewContext(withTokenRequest, withTokenResponse)

// 	testHandler := echo.HandlerFunc(func(context echo.Context) error {
// 		return context.NoContent(http.StatusOK)
// 	})

// 	executableHandler := middle.CSRFmiddle(testHandler)

// 	noTokenErr := executableHandler(noTokenContext)
// 	require.NoError(t, noTokenErr)
// 	assert.Equal(t, noTokenContext.Response().Status, http.StatusForbidden)

// 	withTokenRequest.Header.Set(csrf.CSRFheader, "invalid")
// 	withTokenContext.Set("sid", "test_sid")

// 	withTokenErr := executableHandler(withTokenContext)
// 	require.NoError(t, withTokenErr)
// 	assert.Equal(t, withTokenContext.Response().Status, http.StatusForbidden)

// }

// func TestCheckBoardMemberPermission(t *testing.T) {
// 	// t.Skip()
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	sessionUsecaseMock := sessionMocks.NewMockUseCase(ctrl)
// 	boardUsecaseMock := boardMocks.NewMockUseCase(ctrl)
// 	columnUsecaseMock := columnMocks.NewMockUseCase(ctrl)
// 	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)

// 	middle := middleware.CreateMiddleware(sessionUsecaseMock, boardUsecaseMock, columnUsecaseMock, taskUsecaseMock)

// 	var testUser models.User
// 	err := faker.FakeData(&testUser)
// 	assert.NoError(t, err)
// 	//t.Logf("%+v", testUser)

// 	var testBoard models.Board
// 	err = faker.FakeData(&testBoard)
// 	testBoard.ID = 10
// 	assert.NoError(t, err)
// 	t.Logf("%+v", testBoard)

// 	e := echo.New()

// 	bidRequest := test.NewRequest(echo.PUT, "/", nil)
// 	bidResponse := test.NewRecorder()

// 	bidContext := e.NewContext(bidRequest, bidResponse)
// 	bidContext.SetPath("/boards/:bid")
// 	bidContext.SetParamNames("bid")
// 	//bidContext.SetParamValues(string(testUser.ID))

// 	bidContext.Set("uid", testUser.ID)

// 	// boardUsecaseMock.EXPECT().
// 	// 	Get(testUser.ID, testBoard.ID, false).
// 	// 	Return(nil, nil)

// 	testHandler := echo.HandlerFunc(func(context echo.Context) error {
// 		return context.NoContent(http.StatusOK)
// 	})

// 	executableHandler := middle.CheckBoardAdminPermission(testHandler)

// 	bidErr := executableHandler(bidContext)
// 	//bidErr := testHandler(c)
// 	require.NoError(t, bidErr)
// 	assert.Equal(t, bidContext.Response().Status, http.StatusBadRequest)

// }
