package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	columnMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/mocks"
	columnUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func createRepoMocks(controller *gomock.Controller) *columnMocks.MockRepository {
	columnRepoMock := columnMocks.NewMockRepository(controller)
	return columnRepoMock
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	columnRepoMock := createRepoMocks(ctrl)
	cUsecase := columnUseCase.CreateUseCase(columnRepoMock)

	var testColumn models.Column
	err := faker.FakeData(&testColumn)
	assert.NoError(t, err)
	//t.Logf("%+v", testColumn)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	columnRepoMock.EXPECT().
		Create(&testColumn).
		Return(nil)

	err = cUsecase.Create(&testColumn)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	columnRepoMock := createRepoMocks(ctrl)
	cUsecase := columnUseCase.CreateUseCase(columnRepoMock)

	var testColumn models.Column
	err := faker.FakeData(&testColumn)
	assert.NoError(t, err)
	//t.Logf("%+v", testColumn)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)
	if testBoard.ID == testColumn.Bid {
		testColumn.Bid++
	}
	//t.Logf("%+v", testBoard)

	columnRepoMock.EXPECT().
		Get(testColumn.ID).
		Return(&testColumn, nil)

	column, err := cUsecase.Get(testBoard.ID, testColumn.ID)

	assert.Nil(t, column)
	assert.Equal(t, err, errors.ErrNoPermission)
}

func TestGetTasksByID(t *testing.T) {
	t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	columnRepoMock := createRepoMocks(ctrl)
	cUsecase := columnUseCase.CreateUseCase(columnRepoMock)

	var testColumn models.Column
	err := faker.FakeData(&testColumn)
	assert.NoError(t, err)
	//t.Logf("%+v", testColumn)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	taskSlice := []models.Task{testTask}

	columnRepoMock.EXPECT().
		GetTasksByID(testColumn.ID).
		Return(taskSlice, nil)

	tasks, err := cUsecase.GetTasksByID(testColumn.ID)

	assert.Equal(t, tasks, taskSlice)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	columnRepoMock := createRepoMocks(ctrl)
	cUsecase := columnUseCase.CreateUseCase(columnRepoMock)

	var testColumn models.Column
	err := faker.FakeData(&testColumn)
	assert.NoError(t, err)
	//t.Logf("%+v", testColumn)

	columnRepoMock.EXPECT().
		Delete(testColumn.ID).
		Return(nil)

	err = cUsecase.Delete(testColumn.ID)
	assert.NoError(t, err)
}
