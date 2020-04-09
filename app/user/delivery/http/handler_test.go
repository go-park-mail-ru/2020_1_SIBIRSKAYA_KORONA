package http_test

import (
	"flag"
	"os"
	"testing"

	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/mocks"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/delivery/http"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/mocks"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/usecase"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"net/http"
	test "net/http/httptest"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

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

	os.Exit(m.Run())
}

func createUserHandler(controller *gomock.Controller) *userHandler.UserHandler {
	userRepoMock := userMocks.NewMockRepository(controller)
	sessionRepoMock := sessionMocks.NewMockRepository(controller)

	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	return userHandler.CreateHandlerTest(uUsecase)
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.CreateHandlerTest(userUsecaseMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	body, err := message.GetBody(message.Pair{Name: "user", Data: testUser})

	router := echo.New()

	request := test.NewRequest(echo.GET, "/settings", strings.NewReader(body))
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	userUsecaseMock.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return("test_sid", nil)

	err = handler.Create(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestGet(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.CreateHandlerTest(userUsecaseMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	router := echo.New()

	request := test.NewRequest(echo.GET, "/profile/"+testUser.Nickname, nil)
	t.Log("/profile/" + testUser.Nickname)
	t.Log(request.RequestURI)
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	userUsecaseMock.EXPECT().
		GetByNickname(gomock.Any()).
		Return(&testUser, nil)

	err = handler.Get(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)

	//TODO: сравнить тела в response()

}
