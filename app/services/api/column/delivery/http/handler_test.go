package http_test

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	columnHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/delivery/http"
	columnMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

	test "net/http/httptest"

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

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	columnUsecaseMock := columnMocks.NewMockUseCase(ctrl)
	handler := columnHandler.ColumnHandler{UseCase: columnUsecaseMock}

	var testColumn models.Column
	err := faker.FakeData(&testColumn)
	assert.NoError(t, err)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)

	body, err := testColumn.MarshalJSON()
	assert.NoError(t, err)

	{
		router := echo.New()
		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
		response := test.NewRecorder()
		context := router.NewContext(request, response)
		context.Set("body", body)
		context.Set("bid", testBoard.ID)

		columnUsecaseMock.EXPECT().
			Create(gomock.Any()).
			Return(nil)

		err = handler.Create(context)

		assert.NoError(t, err)
		assert.Equal(t, context.Response().Status, http.StatusOK)
	}

	// {
	// 	router := echo.New()
	// 	request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
	// 	response := test.NewRecorder()
	// 	context := router.NewContext(request, response)
	// 	context.Set("body", body)
	// 	context.Set("tid", testTask.ID)

	// 	checklistUsecaseMock.EXPECT().
	// 		Create(gomock.Any()).
	// 		Return(errors.ErrDbBadOperation)

	// 	err = handler.Create(context)
	// 	if err != nil {
	// 		t.Error(err)
	// 	}
	// 	assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
	// }

}

// func TestGet(t *testing.T) {
// 	// t.Skip()
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	checklistUsecaseMock := checklistMocks.NewMockUseCase(ctrl)
// 	handler := checklistHandler.ChecklistHandler{UseCase: checklistUsecaseMock}

// 	var testChecklist models.Checklist
// 	err := faker.FakeData(&testChecklist)
// 	assert.NoError(t, err)

// 	var testTask models.Task
// 	err = faker.FakeData(&testTask)
// 	assert.NoError(t, err)

// 	{
// 		router := echo.New()
// 		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
// 		response := test.NewRecorder()
// 		context := router.NewContext(request, response)
// 		context.Set("tid", testTask.ID)

// 		checklistUsecaseMock.EXPECT().
// 			Get(testTask.ID).
// 			Return(nil, errors.ErrDbBadOperation)

// 		err = handler.Get(context)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
// 	}

// 	var checklists models.Checklists
// 	checklists = append(checklists, testChecklist)
// 	{
// 		router := echo.New()
// 		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
// 		response := test.NewRecorder()
// 		context := router.NewContext(request, response)
// 		context.Set("tid", testTask.ID)

// 		checklistUsecaseMock.EXPECT().
// 			Get(testTask.ID).
// 			Return(checklists, nil)

// 		err = handler.Get(context)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		assert.Equal(t, context.Response().Status, http.StatusOK)
// 	}

// }
