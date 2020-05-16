package http_test

import (
	"os"
	"testing"

	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	taskHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/delivery/http"
	taskMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/mocks"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	test "net/http/httptest"

	"strings"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)
	handler := taskHandler.TaskHandler{UseCase: taskUsecaseMock}

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)

	body, err := testTask.MarshalJSON()
	assert.NoError(t, err)

	var testColumn models.Column
	err = faker.FakeData(&testColumn)
	assert.NoError(t, err)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	response := test.NewRecorder()
	context := router.NewContext(request, response)
	context.Set("body", body)
	context.Set("cid", testColumn.ID)
	taskUsecaseMock.EXPECT().
		Create(gomock.Any()).
		Return(nil)

	err = handler.Create(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestUpdate(t *testing.T) {
	//t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)
	handler := taskHandler.TaskHandler{UseCase: taskUsecaseMock}

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)

	body, err := testTask.MarshalJSON()
	assert.NoError(t, err)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	context.Set("body", body)
	context.Set("tid", testTask.ID)
	taskUsecaseMock.EXPECT().
		Update(gomock.Any()).
		Return(nil)

	err = handler.Update(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestDelete(t *testing.T) {
	//t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)
	handler := taskHandler.TaskHandler{UseCase: taskUsecaseMock}

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	context.Set("tid", testTask.ID)
	taskUsecaseMock.EXPECT().
		Delete(testTask.ID).
		Return(nil)

	err = handler.Delete(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}
