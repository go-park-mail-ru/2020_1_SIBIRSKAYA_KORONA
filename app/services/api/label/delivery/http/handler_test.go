package http_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	labelHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/delivery/http"
	labelMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/mocks"
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
	labelUsecaseMock := labelMocks.NewMockUseCase(ctrl)
	router := echo.New()
	mw := middleware.CreateMiddleware(nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil)
	labelHandler.CreateHandler(router, labelUsecaseMock, mw)
}

// func TestCreate(t *testing.T) {
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

// 	body, err := testChecklist.MarshalJSON()
// 	assert.NoError(t, err)

// 	{
// 		router := echo.New()
// 		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
// 		response := test.NewRecorder()
// 		context := router.NewContext(request, response)
// 		context.Set("body", body)
// 		context.Set("tid", testTask.ID)

// 		checklistUsecaseMock.EXPECT().
// 			Create(gomock.Any()).
// 			Return(nil)

// 		err = handler.Create(context)

// 		assert.NoError(t, err)
// 		assert.Equal(t, context.Response().Status, http.StatusOK)
// 	}

// 	{
// 		router := echo.New()
// 		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
// 		response := test.NewRecorder()
// 		context := router.NewContext(request, response)
// 		context.Set("body", body)
// 		context.Set("tid", testTask.ID)

// 		checklistUsecaseMock.EXPECT().
// 			Create(gomock.Any()).
// 			Return(errors.ErrDbBadOperation)

// 		err = handler.Create(context)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
// 	}

// }

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

// func TestUpdate(t *testing.T) {
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

// 	body, err := testChecklist.MarshalJSON()
// 	assert.NoError(t, err)

// 	router := echo.New()

// 	request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
// 	response := test.NewRecorder()
// 	context := router.NewContext(request, response)
// 	context.Set("body", body)
// 	context.Set("tid", testTask.ID)

// 	err = handler.Update(context)

// 	assert.NoError(t, err)
// 	assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
// }

// func TestDelete(t *testing.T) {
// 	// t.Skip()
// 	t.Parallel()

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	checklistUsecaseMock := checklistMocks.NewMockUseCase(ctrl)
// 	handler := checklistHandler.ChecklistHandler{UseCase: checklistUsecaseMock}

// 	var testChecklist models.Checklist
// 	err := faker.FakeData(&testChecklist)
// 	assert.NoError(t, err)

// 	{
// 		router := echo.New()
// 		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
// 		response := test.NewRecorder()
// 		context := router.NewContext(request, response)
// 		context.Set("clid", testChecklist.ID)

// 		checklistUsecaseMock.EXPECT().
// 			Delete(gomock.Any()).
// 			Return(nil)

// 		err = handler.Delete(context)
// 		if err != nil {
// 			t.Error(err)
// 		}

// 		assert.NoError(t, err)
// 		assert.Equal(t, context.Response().Status, http.StatusOK)
// 	}

// 	{
// 		router := echo.New()
// 		request := test.NewRequest(echo.POST, "/", strings.NewReader(""))
// 		response := test.NewRecorder()
// 		context := router.NewContext(request, response)
// 		context.Set("clid", testChecklist.ID)

// 		checklistUsecaseMock.EXPECT().
// 			Delete(gomock.Any()).
// 			Return(errors.ErrDbBadOperation)

// 		err = handler.Delete(context)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		assert.Equal(t, context.Response().Status, http.StatusInternalServerError)
// 	}
// }
