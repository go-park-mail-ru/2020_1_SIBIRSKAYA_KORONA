package usecase_test

import (
	"flag"
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	taskMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task/mocks"
	taskUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task/usecase"
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

func createRepoMocks(controller *gomock.Controller) *taskMocks.MockRepository {
	taskRepoMock := taskMocks.NewMockRepository(controller)
	return taskRepoMock
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	taskRepoMock.EXPECT().
		Create(&testTask).
		Return(nil)

	err = tUsecase.Create(&testTask)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	var testColumn models.Column
	err = faker.FakeData(&testColumn)
	assert.NoError(t, err)
	//t.Logf("%+v", testColumn)
	testTask.Cid = testColumn.ID + 1

	taskRepoMock.EXPECT().
		Get(testTask.ID).
		Return(&testTask, nil)

	task, err := tUsecase.Get(testColumn.ID, testTask.ID)
	assert.Nil(t, task)
	assert.Equal(t, err, errors.ErrNoPermission)
}

func TestUpdate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	taskRepoMock.EXPECT().
		Update(testTask).
		Return(nil)

	err = tUsecase.Update(testTask)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	taskRepoMock.EXPECT().
		Delete(testTask.ID).
		Return(nil)

	err = tUsecase.Delete(testTask.ID)
	assert.NoError(t, err)
}
