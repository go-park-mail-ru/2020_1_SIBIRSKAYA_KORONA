package usecase_test

import (
	"flag"
	"os"
	"testing"

	"github.com/bxcodec/faker"
	columnMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column/mocks"
	columnUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
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
	//t.Logf("%+v", testBoard)

	columnRepoMock.EXPECT().
		Get(testColumn.ID).
		Return(&testColumn, nil)

	column, err := cUsecase.Get(testBoard.ID, testColumn.ID)

	assert.Nil(t, column)
	assert.Equal(t, err, errors.ErrNoPermission)
}

func TestGetTasksByID(t *testing.T) {
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
