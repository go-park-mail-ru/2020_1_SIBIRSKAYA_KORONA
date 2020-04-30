package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	itemMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/mocks"
	itemUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createRepoMocks(controller *gomock.Controller) *itemMocks.MockRepository {
	itemRepoMock := itemMocks.NewMockRepository(controller)
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

	itemRepoMock := createRepoMocks(ctrl)
	itUsecase := itemUseCase.CreateUseCase(itemRepoMock)

	var testItem models.Item
	err := faker.FakeData(&testItem)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	itemRepoMock.EXPECT().
		Create(&testItem).
		Return(nil)

	err = itUsecase.Create(&testItem)
	assert.NoError(t, err)

	itemRepoMock.EXPECT().
		Create(&testItem).
		Return(errors.ErrDbBadOperation)

	err = itUsecase.Create(&testItem)
	assert.EqualError(t, err, errors.DbBadOperation)
}

func TestUpdate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepoMock := createRepoMocks(ctrl)
	itUsecase := itemUseCase.CreateUseCase(itemRepoMock)

	var testItem models.Item
	err := faker.FakeData(&testItem)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	itemRepoMock.EXPECT().
		Update(&testItem).
		Return(nil)

	err = itUsecase.Update(&testItem)
	assert.NoError(t, err)

	itemRepoMock.EXPECT().
		Update(&testItem).
		Return(errors.ErrItemNotFound)

	err = itUsecase.Update(&testItem)
	assert.EqualError(t, err, errors.ItemNotFound)
}

func TestGetByID(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepoMock := createRepoMocks(ctrl)
	itUsecase := itemUseCase.CreateUseCase(itemRepoMock)

	var testChecklist models.Checklist
	err := faker.FakeData(&testChecklist)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	var testItem models.Item
	err = faker.FakeData(&testItem)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	testItem.Clid = testChecklist.ID + 1

	itemRepoMock.EXPECT().
		GetByID(testItem.ID).
		Return(&testItem, nil)

	item, err := itUsecase.GetByID(testChecklist.ID, testItem.ID)
	assert.Nil(t, item)
	assert.EqualError(t, err, errors.NoPermission)

	itemRepoMock.EXPECT().
		GetByID(testItem.ID).
		Return(nil, errors.ErrItemNotFound)

	_, err = itUsecase.GetByID(testChecklist.ID, testItem.ID)
	assert.EqualError(t, err, errors.ItemNotFound)

	testItem.Clid = testChecklist.ID
	itemRepoMock.EXPECT().
		GetByID(testItem.ID).
		Return(&testItem, nil)

	item, err = itUsecase.GetByID(testChecklist.ID, testItem.ID)
	assert.NotNil(t, item)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	itemRepoMock := createRepoMocks(ctrl)
	itUsecase := itemUseCase.CreateUseCase(itemRepoMock)

	var testItem models.Item
	err := faker.FakeData(&testItem)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	err = itUsecase.Delete(testItem.ID)
	assert.EqualError(t, err, errors.DbBadOperation)
}
