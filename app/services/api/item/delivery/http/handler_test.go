package http_test

import (
	"os"
	"testing"

	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	itemHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/delivery/http"
	itemMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

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

func TestCreateHandler(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	itemUsecaseMock := itemMocks.NewMockUseCase(ctrl)
	router := echo.New()
	mw := middleware.CreateMiddleware(nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil)
	itemHandler.CreateHandler(router, itemUsecaseMock, mw)
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemUsecaseMock := itemMocks.NewMockUseCase(ctrl)
	handler := itemHandler.ItemHandler{UseCase: itemUsecaseMock}

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)

	var testItem models.Item
	err = faker.FakeData(&testItem)
	assert.NoError(t, err)

	body, err := testItem.MarshalJSON()
	assert.NoError(t, err)

	{
		router := echo.New()
		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
		response := test.NewRecorder()
		context := router.NewContext(request, response)
		context.Set("body", body)
		context.Set("clid", testChecklist.ID)
		itemUsecaseMock.EXPECT().
			Create(gomock.Any()).
			Return(nil)
		err = handler.Create(context)
		assert.NoError(t, err)
		assert.Equal(t, context.Response().Status, http.StatusOK)
	}

	{
		router := echo.New()
		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
		response := test.NewRecorder()
		context := router.NewContext(request, response)
		context.Set("body", body)
		context.Set("clid", testChecklist.ID)
		itemUsecaseMock.EXPECT().
			Create(gomock.Any()).
			Return(errors.ErrDbBadOperation)
		err = handler.Create(context)
		assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
	}
}

func TestUpdate(t *testing.T) {
	//t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemUsecaseMock := itemMocks.NewMockUseCase(ctrl)
	handler := itemHandler.ItemHandler{UseCase: itemUsecaseMock}

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)

	var testItem models.Item
	err = faker.FakeData(&testItem)
	assert.NoError(t, err)

	body, err := testItem.MarshalJSON()
	assert.NoError(t, err)

	{
		router := echo.New()
		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
		response := test.NewRecorder()
		context := router.NewContext(request, response)
		context.Set("body", body)
		context.Set("itid", testItem.ID)
		context.Set("clid", testChecklist.ID)
		itemUsecaseMock.EXPECT().
			Update(gomock.Any()).
			Return(nil)
		err = handler.Update(context)
		assert.NoError(t, err)
		assert.Equal(t, context.Response().Status, http.StatusOK)
	}

	{
		router := echo.New()
		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
		response := test.NewRecorder()
		context := router.NewContext(request, response)
		context.Set("body", body)
		context.Set("itid", testItem.ID)
		context.Set("clid", testChecklist.ID)
		itemUsecaseMock.EXPECT().
			Update(gomock.Any()).
			Return(errors.ErrDbBadOperation)
		err = handler.Update(context)
		assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
	}

}

func TestDelete(t *testing.T) {
	//t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemUsecaseMock := itemMocks.NewMockUseCase(ctrl)
	handler := itemHandler.ItemHandler{UseCase: itemUsecaseMock}

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)

	var testItem models.Item
	err = faker.FakeData(&testItem)
	assert.NoError(t, err)

	{
		router := echo.New()
		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
		response := test.NewRecorder()
		context := router.NewContext(request, response)

		context.Set("itid", testItem.ID)

		err = handler.Delete(context)
		assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
	}
}
