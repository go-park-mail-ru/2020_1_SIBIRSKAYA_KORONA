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
	sessionHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/delivery/http"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/mocks"

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

func TestLogIn(t *testing.T) {
	//t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionUsecaseMock := sessionMocks.NewMockUseCase(ctrl)
	handler := sessionHandler.CreateHandlerTest(sessionUsecaseMock)

	var testUser models.TestUser
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testColumn)

	bodyJSON, err := json.Marshal(testUser)
	body := string(bodyJSON)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/session", strings.NewReader(body))
	response := test.NewRecorder()

	context := router.NewContext(request, response)

	// sessionUsecaseMock.EXPECT().
	// 	Create(gomock.Any(), gomock.Any()).
	// 	Return("test_sid", nil)

	err = handler.LogIn(context)

	assert.NoError(t, err)
	
	assert.Equal(t, context.Response().Status, http.StatusBadRequest)
}