package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	taskMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/mocks"
	taskUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/usecase"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createRepoMocks(controller *gomock.Controller) (*taskMocks.MockRepository, *userMocks.MockRepository) {
	taskRepoMock := taskMocks.NewMockRepository(controller)
	userRepoMock := userMocks.NewMockRepository(controller)
	return taskRepoMock, userRepoMock
}

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock, userRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock, userRepoMock)

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

	taskRepoMock, userRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock, userRepoMock)

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

	taskRepoMock.EXPECT().
		Get(testTask.ID).
		Return(nil, errors.ErrTaskNotFound)

	_, err = tUsecase.Get(testColumn.ID, testTask.ID)
	assert.EqualError(t, err, errors.TaskNotFound)

	taskRepoMock.EXPECT().
		Get(testTask.ID).
		Return(&testTask, nil)

	testTask.Cid = testColumn.ID
	task, err = tUsecase.Get(testColumn.ID, testTask.ID)
	assert.NoError(t, err)
	assert.Equal(t, task, &testTask)
}

func TestUpdate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock, userRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock, userRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	taskRepoMock.EXPECT().
		Update(testTask).
		Return(nil)

	err = tUsecase.Update(testTask)
	assert.NoError(t, err)

	taskRepoMock.EXPECT().
		Update(testTask).
		Return(errors.ErrTaskNotFound)

	err = tUsecase.Update(testTask)
	assert.EqualError(t, err, errors.TaskNotFound)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock, userRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock, userRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	taskRepoMock.EXPECT().
		Delete(testTask.ID).
		Return(nil)

	err = tUsecase.Delete(testTask.ID)
	assert.NoError(t, err)

	taskRepoMock.EXPECT().
		Delete(testTask.ID).
		Return(errors.ErrTaskNotFound)

	err = tUsecase.Delete(testTask.ID)
	assert.EqualError(t, err, errors.TaskNotFound)
}

func TestAssign(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock, userRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock, userRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)
	var testUser models.User
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	userRepoMock.EXPECT().
		GetByID(testUser.ID).
		Return(&testUser, nil)

	taskRepoMock.EXPECT().
		Assign(testTask.ID, &testUser).
		Return(nil)

	err = tUsecase.Assign(testTask.ID, testUser.ID)
	assert.NoError(t, err)

	userRepoMock.EXPECT().
		GetByID(testUser.ID).
		Return(nil, errors.ErrUserNotFound)

	err = tUsecase.Assign(testTask.ID, testUser.ID)
	assert.EqualError(t, err, errors.UserNotFound)

	userRepoMock.EXPECT().
		GetByID(testUser.ID).
		Return(&testUser, nil)

	taskRepoMock.EXPECT().
		Assign(testTask.ID, &testUser).
		Return(errors.ErrDbBadOperation)

	err = tUsecase.Assign(testTask.ID, testUser.ID)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())
}

func TestUnassign(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskRepoMock, userRepoMock := createRepoMocks(ctrl)
	tUsecase := taskUseCase.CreateUseCase(taskRepoMock, userRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)
	var testUser models.User
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	userRepoMock.EXPECT().
		GetByID(testUser.ID).
		Return(&testUser, nil)

	taskRepoMock.EXPECT().
		Unassign(testTask.ID, &testUser).
		Return(nil)

	err = tUsecase.Unassign(testTask.ID, testUser.ID)
	assert.NoError(t, err)

	userRepoMock.EXPECT().
		GetByID(testUser.ID).
		Return(nil, errors.ErrUserNotFound)

	err = tUsecase.Unassign(testTask.ID, testUser.ID)
	assert.EqualError(t, err, errors.UserNotFound)

	userRepoMock.EXPECT().
		GetByID(testUser.ID).
		Return(&testUser, nil)

	taskRepoMock.EXPECT().
		Unassign(testTask.ID, &testUser).
		Return(errors.ErrDbBadOperation)

	err = tUsecase.Unassign(testTask.ID, testUser.ID)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())

}
