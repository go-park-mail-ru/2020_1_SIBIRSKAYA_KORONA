package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	checklistMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist/mocks"
	checklistUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist/usecase"
	itemMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createRepoMocks(controller *gomock.Controller) (*checklistMocks.MockRepository, *itemMocks.MockRepository) {
	checklistRepoMock := checklistMocks.NewMockRepository(controller)
	itemRepoMock := itemMocks.NewMockRepository(controller)
	return checklistRepoMock, itemRepoMock
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

	checklistRepoMock, itemRepoMock := createRepoMocks(ctrl)
	chUsecase := checklistUseCase.CreateUseCase(checklistRepoMock, itemRepoMock)

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	checklistRepoMock.EXPECT().
		Create(&testChecklist).
		Return(nil)

	err = chUsecase.Create(&testChecklist)
	assert.NoError(t, err)

	checklistRepoMock.EXPECT().
		Create(&testChecklist).
		Return(errors.ErrDbBadOperation)

	err = chUsecase.Create(&testChecklist)
	assert.EqualError(t, err, errors.DbBadOperation)

}

func TestGet(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	checklistRepoMock, itemRepoMock := createRepoMocks(ctrl)
	chUsecase := checklistUseCase.CreateUseCase(checklistRepoMock, itemRepoMock)

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	var testChecklists []models.Checklist
	testChecklists = append(testChecklists, testChecklist)

	checklistRepoMock.EXPECT().
		Get(testTask.ID).
		Return(testChecklists, nil)

	checklists, err := chUsecase.Get(testTask.ID)
	assert.NotNil(t, checklists)
	assert.NoError(t, err)

	checklistRepoMock.EXPECT().
		Get(testTask.ID).
		Return(nil, errors.ErrChecklistNotFound)

	_, err = chUsecase.Get(testTask.ID)
	assert.EqualError(t, err, errors.ChecklistNotFound)
}

func TestGetByID(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	checklistRepoMock, itemRepoMock := createRepoMocks(ctrl)
	chUsecase := checklistUseCase.CreateUseCase(checklistRepoMock, itemRepoMock)

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	testChecklist.Tid = testTask.ID + 1

	checklistRepoMock.EXPECT().
		GetByID(testChecklist.ID).
		Return(&testChecklist, nil)

	checklists, err := chUsecase.GetByID(testTask.ID, testChecklist.ID)
	assert.Nil(t, checklists)
	assert.EqualError(t, err, errors.NoPermission)

	checklistRepoMock.EXPECT().
		GetByID(testChecklist.ID).
		Return(nil, errors.ErrChecklistNotFound)
	_, err = chUsecase.GetByID(testTask.ID, testChecklist.ID)
	assert.EqualError(t, err, errors.ChecklistNotFound)

	testChecklist.Tid = testTask.ID
	checklistRepoMock.EXPECT().
		GetByID(testChecklist.ID).
		Return(&testChecklist, nil)

	checklist, err := chUsecase.GetByID(testTask.ID, testChecklist.ID)
	assert.NotNil(t, checklist)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	checklistRepoMock, itemRepoMock := createRepoMocks(ctrl)
	chUsecase := checklistUseCase.CreateUseCase(checklistRepoMock, itemRepoMock)

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	checklistRepoMock.EXPECT().
		Delete(testChecklist.ID).
		Return(nil)

	err = chUsecase.Delete(testChecklist.ID)
	assert.NoError(t, err)

	checklistRepoMock.EXPECT().
		Delete(testChecklist.ID).
		Return(errors.ErrChecklistNotFound)

	err = chUsecase.Delete(testChecklist.ID)
	assert.EqualError(t, err, errors.ChecklistNotFound)
}
