package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	labelMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/mocks"
	labelUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/usecase"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createRepoMocks(controller *gomock.Controller) *labelMocks.MockRepository {
	itemRepoMock := labelMocks.NewMockRepository(controller)
	return itemRepoMock
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

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		Create(&testLabel).
		Return(nil)

	err = lUsecase.Create(&testLabel)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		Create(&testLabel).
		Return(errors.ErrDbBadOperation)

	err = lUsecase.Create(&testLabel)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())
}

func TestGet(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)

	testLabel.Bid = testBoard.ID + 1

	labelRepoMock.EXPECT().
		Get(testLabel.ID).
		Return(&testLabel, nil)

	label, err := lUsecase.Get(testBoard.ID, testLabel.ID)

	assert.Nil(t, label)
	assert.EqualError(t, err, errors.ErrNoPermission.Error())

	testLabel.Bid = testBoard.ID

	labelRepoMock.EXPECT().
		Get(testLabel.ID).
		Return(&testLabel, nil)

	label, err = lUsecase.Get(testBoard.ID, testLabel.ID)

	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		Get(testLabel.ID).
		Return(nil, errors.ErrLabelNotFound)

	label, err = lUsecase.Get(testBoard.ID, testLabel.ID)
	assert.Nil(t, label)
	assert.EqualError(t, err, errors.ErrLabelNotFound.Error())
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	labelRepoMock.EXPECT().
		Delete(testLabel.ID).
		Return(nil)

	err = lUsecase.Delete(testLabel.ID)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		Delete(testLabel.ID).
		Return(errors.ErrLabelNotFound)

	err = lUsecase.Delete(testLabel.ID)
	assert.EqualError(t, err, errors.ErrLabelNotFound.Error())
}

func TestUpdate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	labelRepoMock.EXPECT().
		Update(testLabel).
		Return(nil)

	err = lUsecase.Update(testLabel)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		Update(testLabel).
		Return(errors.ErrLabelNotFound)

	err = lUsecase.Update(testLabel)
	assert.EqualError(t, err, errors.ErrLabelNotFound.Error())
}

func TestAddLabelOnTask(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		AddLabelOnTask(testLabel.ID, testTask.ID).
		Return(nil)

	err = lUsecase.AddLabelOnTask(testLabel.ID, testTask.ID)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		AddLabelOnTask(testLabel.ID, testTask.ID).
		Return(errors.ErrTaskNotFound)

	err = lUsecase.AddLabelOnTask(testLabel.ID, testTask.ID)
	assert.EqualError(t, err, errors.ErrTaskNotFound.Error())
}

func TestRemoveLabelFromTask(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		RemoveLabelFromTask(testLabel.ID, testTask.ID).
		Return(nil)

	err = lUsecase.RemoveLabelFromTask(testLabel.ID, testTask.ID)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		RemoveLabelFromTask(testLabel.ID, testTask.ID).
		Return(errors.ErrLabelNotFound)

	err = lUsecase.RemoveLabelFromTask(testLabel.ID, testTask.ID)
	assert.EqualError(t, err, errors.ErrLabelNotFound.Error())
}
