package usecase_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/mocks"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/usecase"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func createRepoMocks(controller *gomock.Controller) (*userMocks.MockRepository, *sessionMocks.MockRepository) {
	userRepoMock := userMocks.NewMockRepository(controller)
	sessionRepoMock := sessionMocks.NewMockRepository(controller)
	return userRepoMock, sessionRepoMock
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	ses := &models.Session{ID: testUser.ID, Expires: 5}

	// good
	sid := "test_sid"
	userRepoMock.EXPECT().GetByNickname(testUser.Nickname).Return(&testUser, nil)
	userRepoMock.EXPECT().CheckPassword(testUser.ID, testUser.Password).Return(true)
	sessionRepoMock.EXPECT().Create(ses).Return(sid, nil)
	resSid, err := sUsecase.Create(&testUser, ses.Expires)
	assert.NoError(t, err)
	assert.Equal(t, sid, resSid)
}

/*
func TestGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	testSid := "test_sid"
	sessionRepoMock.EXPECT().
		Get("test_sid").
		Return(uint(10), true)

	uid, exist := sUsecase.Get(testSid)

	assert.Equal(t, uid, uint(10))
	assert.True(t, exist)

}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	testSid := "test_sid"
	sessionRepoMock.EXPECT().
		Delete("test_sid").
		Return(nil)

	err := sUsecase.Delete(testSid)

	assert.NoError(t, err)
}
*/
