package usecase_test

import (
	"os"
	"testing"

	"mime/multipart"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	attachMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/mocks"
	attachUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createRepoMocks(controller *gomock.Controller) (*attachMocks.MockRepository, *attachMocks.MockFileRepository) {
	itemRepoMock := attachMocks.NewMockRepository(controller)
	itemFileRepoMock := attachMocks.NewMockFileRepository(controller)
	return itemRepoMock, itemFileRepoMock
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

	attachRepoMock, attachFileRepoMock := createRepoMocks(ctrl)
	atchUsecase := attachUseCase.CreateUseCase(attachRepoMock, attachFileRepoMock)

	var testAttach models.AttachedFile
	err := faker.FakeData(&testAttach)
	assert.NoError(t, err)

	fakeAttach := multipart.FileHeader{Filename: "fake"}

	attachFileRepoMock.EXPECT().
		UploadFile(&fakeAttach, &testAttach).
		Return(testAttach.URL, nil)

	attachRepoMock.EXPECT().
		Create(&testAttach).
		Return(nil)

	err = atchUsecase.Create(&testAttach, &fakeAttach)
	assert.NoError(t, err)

	attachFileRepoMock.EXPECT().
		UploadFile(&fakeAttach, &testAttach).
		Return("", errors.ErrBadFileUploadS3)

	err = atchUsecase.Create(&testAttach, &fakeAttach)
	assert.EqualError(t, err, errors.BadFileUploadS3)

	attachFileRepoMock.EXPECT().
		UploadFile(&fakeAttach, &testAttach).
		Return(testAttach.URL, nil)

	attachRepoMock.EXPECT().
		Create(&testAttach).
		Return(errors.ErrDbBadOperation)

	err = atchUsecase.Create(&testAttach, &fakeAttach)
	assert.EqualError(t, err, errors.DbBadOperation)
}

func TestDelete(t *testing.T) {
	t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	attachRepoMock, attachFileRepoMock := createRepoMocks(ctrl)
	atchUsecase := attachUseCase.CreateUseCase(attachRepoMock, attachFileRepoMock)

	var testAttach models.AttachedFile
	err := faker.FakeData(&testAttach)
	assert.NoError(t, err)

	attachRepoMock.EXPECT().
		GetByID(testAttach.ID).
		Return(&testAttach, nil)

	attachFileRepoMock.EXPECT().
		DeleteFile(testAttach.Name).
		Return(nil)

	attachRepoMock.EXPECT().
		Delete(testAttach.ID).
		Return(nil)

	err = atchUsecase.Delete(testAttach.ID)
	assert.NoError(t, err)

	attachRepoMock.EXPECT().
		GetByID(testAttach.ID).
		Return(nil, errors.ErrDbBadOperation)

	err = atchUsecase.Delete(testAttach.ID)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())

	attachRepoMock.EXPECT().
		GetByID(testAttach.ID).
		Return(&testAttach, nil)

	attachFileRepoMock.EXPECT().
		DeleteFile(testAttach.Name).
		Return(errors.ErrDbBadOperation)

	err = atchUsecase.Delete(testAttach.ID)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())

	attachRepoMock.EXPECT().
		GetByID(testAttach.ID).
		Return(&testAttach, nil)

	attachFileRepoMock.EXPECT().
		DeleteFile(testAttach.Name).
		Return(nil)

	attachRepoMock.EXPECT().
		Delete(testAttach.ID).
		Return(errors.ErrDbBadOperation)

	err = atchUsecase.Delete(testAttach.ID)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())
}

func TestGet(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	attachRepoMock, attachFileRepoMock := createRepoMocks(ctrl)
	atchUsecase := attachUseCase.CreateUseCase(attachRepoMock, attachFileRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)

	var testAttach models.AttachedFile
	err = faker.FakeData(&testAttach)
	assert.NoError(t, err)

	var attachs []models.AttachedFile
	attachs = append(attachs, testAttach)

	attachRepoMock.EXPECT().
		Get(testTask.ID).
		Return(attachs, nil)

	getAttachs, err := atchUsecase.Get(testTask.ID)
	assert.NotNil(t, getAttachs)
	assert.NoError(t, err)

	attachRepoMock.EXPECT().
		Get(testTask.ID).
		Return(nil, errors.ErrDbBadOperation)

	getAttachs, err = atchUsecase.Get(testTask.ID)
	assert.Nil(t, getAttachs)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())
}

func TestGetByID(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	attachRepoMock, attachFileRepoMock := createRepoMocks(ctrl)
	atchUsecase := attachUseCase.CreateUseCase(attachRepoMock, attachFileRepoMock)

	var testTask models.Task
	err := faker.FakeData(&testTask)
	assert.NoError(t, err)

	var testAttach models.AttachedFile
	err = faker.FakeData(&testAttach)
	assert.NoError(t, err)

	testAttach.Tid = testTask.ID + 1

	attachRepoMock.EXPECT().
		GetByID(testAttach.ID).
		Return(&testAttach, nil)

	attach, err := atchUsecase.GetByID(testTask.ID, testAttach.ID)
	assert.Nil(t, attach)
	assert.EqualError(t, err, errors.NoPermission)

	testAttach.Tid = testTask.ID
	attachRepoMock.EXPECT().
		GetByID(testAttach.ID).
		Return(&testAttach, nil)

	attach, err = atchUsecase.GetByID(testTask.ID, testAttach.ID)
	assert.NotNil(t, attach)
	assert.NoError(t, err)

	attachRepoMock.EXPECT().
		GetByID(testAttach.ID).
		Return(nil, errors.ErrDbBadOperation)

	attach, err = atchUsecase.GetByID(testTask.ID, testAttach.ID)
	assert.Nil(t, attach)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())
}
