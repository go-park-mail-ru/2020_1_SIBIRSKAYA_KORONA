package middleware_test

import (
	"flag"
	"os"
	"testing"

	"net/http"
	test "net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
)

// TODO: поднять отдельный пакет, в котором будет общие параметры
var test_opts struct {
	configPath string
}

// должны поднять конфиг для тестов
func TestMain(m *testing.M) {
	flag.StringVar(&test_opts.configPath, "test-c", "", "path to configuration file")
	flag.StringVar(&test_opts.configPath, "test-config", "", "path to configuration file")
	flag.Parse()

	viper.SetConfigFile(test_opts.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestCORS(t *testing.T) {
	e := echo.New()
	request := test.NewRequest(echo.GET, "/", nil)
	response := test.NewRecorder()
	context := e.NewContext(request, response)
	middle := middleware.InitMiddleware()

	executableHandler := middle.CORS(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))

	err := executableHandler(context)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:5757", response.Header().Get("Access-Control-Allow-Origin"))
}

func TestPanicProcess(t *testing.T) {
	e := echo.New()
	request := test.NewRequest(echo.GET, "/settings", nil)
	response := test.NewRecorder()
	context := e.NewContext(request, response)
	middle := middleware.InitMiddleware()

	panicHandler := echo.HandlerFunc(func(c echo.Context) error {
		if 2+2 == 4 {
			panic("big panic")
		}
		return c.NoContent(http.StatusOK)
	})

	processedPanicHandler := middle.ProcessPanic(panicHandler)

	err := processedPanicHandler(context)
	require.NoError(t, err)
}

func TestCheckCookieExist(t *testing.T) {
	e := echo.New()

	middle := middleware.InitMiddleware()

	noCookieRequest := test.NewRequest(echo.GET, "/", nil)
	noCookieResponse := test.NewRecorder()
	noCookieContext := e.NewContext(noCookieRequest, noCookieResponse)

	withCookieRequest := test.NewRequest(echo.GET, "/", nil)
	withCookieResponse := test.NewRecorder()
	withCookieContext := e.NewContext(withCookieRequest, withCookieResponse)

	cookie := http.Cookie{Name: "session_id", Value: "check_only_for_exist"}
	withCookieContext.Request().AddCookie(&cookie)

	testHandler := echo.HandlerFunc(func(context echo.Context) error {
		return context.NoContent(http.StatusOK)
	})

	executableHandler := middle.CheckCookieExist(testHandler)

	noCookieErr := executableHandler(noCookieContext)
	require.NoError(t, noCookieErr)
	assert.Equal(t, noCookieContext.Response().Status, http.StatusForbidden)

	withCookieErr := executableHandler(withCookieContext)
	require.NoError(t, withCookieErr)
	assert.Equal(t, withCookieContext.Response().Status, http.StatusOK)
}
