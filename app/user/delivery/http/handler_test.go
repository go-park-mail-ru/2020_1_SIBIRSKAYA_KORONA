package http_test

import (
	"flag"
	"io"
	"os"
	"testing"

	"encoding/json"

	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/delivery/http"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/mocks"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/usecase"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	test "net/http/httptest"

	"path/filepath"
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

	var testUser models.TestUser
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	bodyJSON, err := json.Marshal(testUser)
	body := string(bodyJSON)

	router := echo.New()

	request := test.NewRequest(echo.POST, "/settings", strings.NewReader(body))
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	// userUsecaseMock.EXPECT().
	// 	Create(gomock.Any(), gomock.Any()).
	// 	Return("test_sid", nil)

	err = handler.Create(context)

	assert.NoError(t, err)
	//assert.Equal(t, context.Response().Status, http.StatusOK)
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
	//t.Log("/profile/" + testUser.Nickname)
	response := test.NewRecorder()
	context := router.NewContext(request, response)

	userUsecaseMock.EXPECT().
		GetByNickname(gomock.Any()).
		Return(&testUser, nil)

	err = handler.Get(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)

}

func TestGetAll(t *testing.T) {
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

	response := test.NewRecorder()
	context := router.NewContext(request, response)

	context.Set("uid", testUser.ID)

	userUsecaseMock.EXPECT().
		GetByID(testUser.ID).
		Return(&testUser, nil)

	err = handler.GetAll(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestGetBoards(t *testing.T) {
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

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	boardSlice := []models.Board{testBoard}

	router := echo.New()
	request := test.NewRequest(echo.GET, "/profile/"+testUser.Nickname, nil)

	response := test.NewRecorder()
	context := router.NewContext(request, response)

	context.Set("uid", testUser.ID)

	userUsecaseMock.EXPECT().
		GetBoardsByID(testUser.ID).
		Return(boardSlice, boardSlice, nil)

	err = handler.GetBoards(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}

func TestGetUpdate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUsecaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.CreateHandlerTest(userUsecaseMock)

	var testUser models.TestUser
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	pathToDefaultAvatar := fmt.Sprintf("%s/%s",
		viper.GetString("frontend.public_dir"),
		viper.GetString("frontend.default_avatar"))

	file, err := os.Open(pathToDefaultAvatar)
	assert.NoError(t, err)
	defer file.Close()

	router := echo.New()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("newName", testUser.Name)
	writer.WriteField("newSurname", testUser.Surname)
	writer.WriteField("newNickname", testUser.Nickname)
	writer.WriteField("newEmail", testUser.Email)
	writer.WriteField("newPassword", testUser.Password+"postfix")
	writer.WriteField("oldPassword", testUser.Password)

	part, err := writer.CreateFormFile("avatar", filepath.Base(pathToDefaultAvatar))
	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	request := test.NewRequest(echo.PUT, "/settings", body)

	response := test.NewRecorder()
	context := router.NewContext(request, response)
	context.Set("uid", testUser.ID)

	userUsecaseMock.EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
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

	userUsecaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.CreateHandlerTest(userUsecaseMock)

	var testUser models.TestUser
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	router := echo.New()
	request := test.NewRequest(echo.DELETE, "/settings", nil)

	response := test.NewRecorder()
	context := router.NewContext(request, response)
	context.Set("uid", testUser.ID)
	context.Set("sid", "e24f326c-c489-4c5f-9cea-4b9681877155")

	userUsecaseMock.EXPECT().
		Delete(testUser.ID, "e24f326c-c489-4c5f-9cea-4b9681877155").
		Return(nil)

	err = handler.Delete(context)

	assert.NoError(t, err)
	assert.Equal(t, context.Response().Status, http.StatusOK)
}
