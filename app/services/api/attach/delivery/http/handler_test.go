package http_test

import (
	"bytes"
	"mime/multipart"
	"os"
	"strings"
	"testing"

	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	attachHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/delivery/http"
	attachMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/mocks"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	test "net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	attachUsecaseMock := attachMocks.NewMockUseCase(ctrl)
	handler := attachHandler.AttachHandler{UseCase: attachUsecaseMock}

	var testAttach models.AttachedFile
	err := faker.FakeData(&testAttach)
	assert.NoError(t, err)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)

	router := echo.New()

	filebody := &bytes.Buffer{}
	writer := multipart.NewWriter(filebody)
	writer.CreateFormFile("file", "")
	writer.WriteField("test", "test")
	err = writer.Close()

	request := test.NewRequest(echo.POST, "/", filebody)
	response := test.NewRecorder()
	context := router.NewContext(request, response)
	context.Set("tid", testTask.ID)

	attachUsecaseMock.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	err = handler.Create(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	attachUsecaseMock := attachMocks.NewMockUseCase(ctrl)
	handler := attachHandler.AttachHandler{UseCase: attachUsecaseMock}

	var testAttach models.AttachedFile
	err := faker.FakeData(&testAttach)
	assert.NoError(t, err)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	response := test.NewRecorder()
	context := router.NewContext(request, response)
	context.Set("fid", testAttach.ID)

	attachUsecaseMock.EXPECT().
		Delete(gomock.Any()).
		Return(nil)

	err = handler.Delete(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}
