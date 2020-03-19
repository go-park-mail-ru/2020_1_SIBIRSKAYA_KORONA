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
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewRecorder()
	c := e.NewContext(req, res)
	m := middleware.InitMiddleware()

	h := m.CORS(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))

	err := h(c)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:5757", res.Header().Get("Access-Control-Allow-Origin"))
}

func TestPanicProcess(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewRecorder()
	c := e.NewContext(req, res)
	m := middleware.InitMiddleware()

	panicHandler := echo.HandlerFunc(func(c echo.Context) error {
		if 2+2 == 4 {
			panic("oamoamaoama")
		}
		return c.NoContent(http.StatusOK)
	})

	processedPanicHandler := m.ProcessPanic(panicHandler)

	err := processedPanicHandler(c)
	require.NoError(t, err)
}
