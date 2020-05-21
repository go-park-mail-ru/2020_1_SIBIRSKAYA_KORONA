package usecase_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	ntfMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/mocks"
	ntfUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/usecase"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func createRepoMocks(controller *gomock.Controller) (*userMocks.MockRepository, *ntfMocks.MockRepository) {
	usrRepoMock := userMocks.NewMockRepository(controller)
	ntfRepoMock := ntfMocks.NewMockRepository(controller)
	return usrRepoMock, ntfRepoMock
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usrRepoMock, ntfRepoMock := createRepoMocks(ctrl)
	ntfUsecase := ntfUseCase.CreateUseCase(usrRepoMock, ntfRepoMock)

	var testEvent models.Event
	err := faker.FakeData(&testEvent)
	assert.NoError(t, err)

	// good
	ntfRepoMock.EXPECT().Create(&testEvent).Return(nil)
	err = ntfUsecase.Create(&testEvent)
	assert.NoError(t, err)

	// err
	ntfRepoMock.EXPECT().Create(&testEvent).Return(errors.ErrConflict)
	err = ntfUsecase.Create(&testEvent)
	assert.Equal(t, err, errors.ErrConflict)
}

func TestGetAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usrRepoMock, ntfRepoMock := createRepoMocks(ctrl)
	ntfUsecase := ntfUseCase.CreateUseCase(usrRepoMock, ntfRepoMock)

	var testEvents models.Events
	err := faker.FakeData(&testEvents)
	assert.NoError(t, err)

	// good
	var uid uint = 10
	ntfRepoMock.EXPECT().GetAll(uid).Return(testEvents, true)
	for idx := range testEvents {
		usrRepoMock.EXPECT().GetByID(testEvents[idx].MakeUid).Return(testEvents[idx].MakeUsr, nil)
		if testEvents[idx].MetaData.Uid != 0 {
			usrRepoMock.EXPECT().GetByID(testEvents[idx].MetaData.Uid).Return(testEvents[idx].MetaData.Usr, nil)
		}
	}
	res, has := ntfUsecase.GetAll(uid)
	assert.Equal(t, has, true)
	assert.Equal(t, res, testEvents)

	// error
	uid++
	ntfRepoMock.EXPECT().GetAll(uid).Return(nil, false)
	_, has = ntfUsecase.GetAll(uid)
	assert.Equal(t, has, false)
}

func TestUpdateAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usrRepoMock, ntfRepoMock := createRepoMocks(ctrl)
	ntfUsecase := ntfUseCase.CreateUseCase(usrRepoMock, ntfRepoMock)

	// good
	var uid uint = 10
	ntfRepoMock.EXPECT().UpdateAll(uid).Return(nil)
	err := ntfUsecase.UpdateAll(uid)
	assert.NoError(t, err)

	// err
	uid++
	ntfRepoMock.EXPECT().UpdateAll(uid).Return(errors.ErrDbBadOperation)
	err = ntfUsecase.UpdateAll(uid)
	assert.Equal(t, err, errors.ErrDbBadOperation)
}

func TestDeleteAll(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usrRepoMock, ntfRepoMock := createRepoMocks(ctrl)
	ntfUsecase := ntfUseCase.CreateUseCase(usrRepoMock, ntfRepoMock)

	// good
	var uid uint = 10
	ntfRepoMock.EXPECT().DeleteAll(uid).Return(nil)
	err := ntfUsecase.DeleteAll(uid)
	assert.NoError(t, err)

	// err
	uid++
	ntfRepoMock.EXPECT().DeleteAll(uid).Return(errors.ErrDbBadOperation)
	err = ntfUsecase.DeleteAll(uid)
	assert.Equal(t, err, errors.ErrDbBadOperation)
}