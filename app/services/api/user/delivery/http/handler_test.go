package http_test

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	test "net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/labstack/echo/v4"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/delivery/http"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func GetContexFromJSON(method, path string) echo.Context {
	request := test.NewRequest(method, path, nil)
	return echo.New().NewContext(request, test.NewRecorder())
}

/*func GetContexFromMultiPart(method, path string, body *bytes.Buffer) echo.Context {
	request := test.NewRequest(method, path, body)
	return echo.New().NewContext(request, test.NewRecorder())
}*/

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	// good
	{
		body, err := testUser.MarshalJSON()
		assert.NoError(t, err)
		userUseCaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return("test_sid", nil)
		ctx := GetContexFromJSON(echo.POST, "/settings")
		ctx.Set("body", body)
		err = handler.Create(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}
	// error
	{
		testUser.ID++
		body, err := testUser.MarshalJSON()
		assert.NoError(t, err)
		userUseCaseMock.EXPECT().Create(gomock.Any(), gomock.Any()).Return("", errors.ErrConflict)
		ctx := GetContexFromJSON(echo.POST, "/settings")
		ctx.Set("body", body)
		err = handler.Create(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusConflict)
	}

}

func TestGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	{
		ctx := GetContexFromJSON(echo.GET, "/profile/"+testUser.Nickname)
		ctx.SetParamNames("id_or_nickname")

		// good nickname
		userUseCaseMock.EXPECT().GetByNickname(testUser.Nickname).Return(&testUser, nil)
		ctx.SetParamValues(testUser.Nickname)
		err = handler.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	{
		strID := strconv.Itoa(int(testUser.ID))
		ctx := GetContexFromJSON(echo.GET, "/profile/"+strID)
		ctx.SetParamNames("id_or_nickname")

		// good id
		userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
		ctx.SetParamValues(strID)
		err = handler.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	{
		testUser.Nickname += "aa"
		ctx := GetContexFromJSON(echo.GET, "/profile/"+testUser.Nickname)
		ctx.SetParamNames("id_or_nickname")

		// error nickname
		userUseCaseMock.EXPECT().GetByNickname(testUser.Nickname).Return(nil, errors.ErrUserNotFound)
		ctx.SetParamValues(testUser.Nickname)
		err = handler.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusNotFound)
	}

	{
		testUser.ID++
		strID := strconv.Itoa(int(testUser.ID))
		ctx := GetContexFromJSON(echo.GET, "/profile/"+strID)
		ctx.SetParamNames("id_or_nickname")

		// error id
		userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
		ctx.SetParamValues(strID)
		err = handler.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusNotFound)
	}
}

func TestGetAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	{
		ctx := GetContexFromJSON(echo.GET, "/settings")
		ctx.Set("uid", testUser.ID)

		// good
		userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
		err = handler.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	{
		ctx := GetContexFromJSON(echo.GET, "/settings")
		testUser.ID++
		ctx.Set("uid", testUser.ID)

		// error
		userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
		err = handler.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusNotFound)
	}
}

func TestGetUpdate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	testUser.Password = []byte("bbbb")
	oldPass := string(testUser.Password)
	newPass := string(testUser.Password) + "aaa"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err = writer.WriteField("newName", testUser.Name)
	if err != nil {
		t.Error(err)
	}
	err = writer.WriteField("newSurname", testUser.Surname)
	if err != nil {
		t.Error(err)
	}
	err = writer.WriteField("newNickname", testUser.Nickname)
	if err != nil {
		t.Error(err)
	}
	err = writer.WriteField("newEmail", testUser.Email)
	if err != nil {
		t.Error(err)
	}
	err = writer.WriteField("newPassword", newPass)
	if err != nil {
		t.Error(err)
	}
	err = writer.WriteField("oldPassword", oldPass)
	if err != nil {
		t.Error(err)
	}
	err = writer.Close()
	assert.NoError(t, err)
	{
		request := test.NewRequest(echo.PUT, "/settings", body)
		ctx := echo.New().NewContext(request, test.NewRecorder())
		//ctx := GetContexFromMultiPart(echo.PUT, "/settings", body)
		ctx.Set("uid", testUser.ID)
		// good
		testUser.Password = []byte(newPass)
		// TODO: разобраться без any
		userUseCaseMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		err = handler.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}
	{
		request := test.NewRequest(echo.PUT, "/settings", body)
		ctx := echo.New().NewContext(request, test.NewRecorder())
		testUser.ID++
		ctx.Set("uid", testUser.ID)
		// error
		userUseCaseMock.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.ErrUserNotFound)
		err = handler.Update(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusNotFound)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}

	var id uint = 1
	sid := "test_sid"

	{
		ctx := GetContexFromJSON(echo.DELETE, "/settings")
		ctx.Set("uid", id)
		ctx.Set("sid", sid)

		// good
		userUseCaseMock.EXPECT().Delete(id, sid).Return(nil)
		err := handler.Delete(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	{
		id++

		ctx := GetContexFromJSON(echo.DELETE, "/settings")
		ctx.Set("uid", id)
		ctx.Set("sid", sid)

		// error
		userUseCaseMock.EXPECT().Delete(id, sid).Return(errors.ErrUserNotFound)
		err := handler.Delete(ctx)
		assert.NoError(t, err)
		log.Println(ctx.Response().Status, http.StatusNotFound)
		assert.Equal(t, ctx.Response().Status, http.StatusNotFound)
	}
}

func TestGetUsersByNicknamePart(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCaseMock := userMocks.NewMockUseCase(ctrl)
	handler := userHandler.UserHandler{UseCase: userUseCaseMock}

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	{
		ctx := GetContexFromJSON(echo.GET, "/settings")
		ctx.Set("uid", testUser.ID)

		// good
		userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
		err = handler.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusOK)
	}

	{
		ctx := GetContexFromJSON(echo.GET, "/settings")
		testUser.ID++
		ctx.Set("uid", testUser.ID)

		// error
		userUseCaseMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
		err = handler.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, http.StatusNotFound)
	}
}

/*

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
/*
*/
