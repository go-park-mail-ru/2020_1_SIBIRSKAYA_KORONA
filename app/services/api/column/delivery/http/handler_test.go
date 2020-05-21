package http_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	columnHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/delivery/http"
	columnMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func TestCreateHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	columnUsecaseMock := columnMocks.NewMockUseCase(ctrl)
	router := echo.New()
	mw := middleware.CreateMiddleware(nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil)
	columnHandler.CreateHandler(router, columnUsecaseMock, mw)
}
