package http_test

import (
	"flag"
	"os"
	"testing"

	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	taskHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task/delivery/http"
	taskMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task/mocks"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	test "net/http/httptest"

	"strings"

	"github.com/stretchr/testify/assert"
)

var test_opts struct {
	configPath string
}

func TestMain(m *testing.M) {
	flag.StringVar(&test_opts.configPath, "test-c", "", "path to configuration file")
	flag.StringVar(&test_opts.configPath, "test-config", "", "path to configuration file")
	flag.Parse()

	viper.SetConfigFile(test_opts.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	logger.InitLogger()

	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)
	handler := taskHandler.CreateHandlerTest(taskUsecaseMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	var testColumn models.Column
	err = faker.FakeData(&testColumn)
	assert.NoError(t, err)
	//t.Logf("%+v", testColumn)

	bodyJSON, err := json.Marshal(testTask)
	body := string(bodyJSON)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/", strings.NewReader(body))
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	context.Set("cid", testColumn.ID)
	taskUsecaseMock.EXPECT().
		Create(gomock.Any()).
		Return(nil)

	err = handler.Create(context)

	assert.NoError(t, err)
	//assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestUpdate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)
	handler := taskHandler.CreateHandlerTest(taskUsecaseMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	bodyJSON, err := json.Marshal(testTask)
	body := string(bodyJSON)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/", strings.NewReader(body))
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	context.Set("tid", testTask.ID)
	taskUsecaseMock.EXPECT().
		Update(gomock.Any()).
		Return(nil)

	err = handler.Update(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskUsecaseMock := taskMocks.NewMockUseCase(ctrl)
	handler := taskHandler.CreateHandlerTest(taskUsecaseMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	bodyJSON, err := json.Marshal(testTask)
	body := string(bodyJSON)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/", strings.NewReader(body))
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
