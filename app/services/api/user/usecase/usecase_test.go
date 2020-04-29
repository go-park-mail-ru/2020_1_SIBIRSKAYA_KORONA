package usecase_test

import (
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/mocks"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/usecase"
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

	userRepoMock.EXPECT().
		Create(gomock.Any()).
		Return(nil)
	sessionRepoMock.EXPECT().
		Create(gomock.Any()).
		Return("cookie_value") //, nil)

	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	sessionExpires := time.Now().AddDate(1, 0, 0)

	sid, err := uUsecase.Create(&testUser, sessionExpires)

	assert.NoError(t, err)
	assert.Equal(t, sid, "cookie_value")
}

func TestGetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)

	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	ID := testUser.ID

	userRepoMock.EXPECT().
		GetByID(ID).
		Return(&testUser, nil)

	user, err := uUsecase.GetByID(ID)
	assert.NoError(t, err)
	assert.Equal(t, &testUser, user)

}

func TestGetByNickname(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)

	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	nickname := testUser.Nickname

	userRepoMock.EXPECT().
		GetByNickname(nickname).
		Return(&testUser, nil)

	user, err := uUsecase.GetByNickname(nickname)
	assert.NoError(t, err)
	assert.Equal(t, &testUser, user)
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)

	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	oldPass := []byte("oldPass")

	// передать nil в качестве *multipart.FileHeader приемлемо, это самый частый кейс
	// работу с картинкой будем гонять на тестах юзерского репозитория
	userRepoMock.EXPECT().
		Update([]byte("oldPass"), &testUser, nil).
		Return(nil)

	updateErr := uUsecase.Update(oldPass, &testUser, nil)

	assert.NoError(t, updateErr)

	// второго вызова к юзерской репе произойти не должно
	updateErr = uUsecase.Update([]byte("no"), nil, nil)

	assert.Equal(t, updateErr, errors.ErrInternal)

}

func TestGetBoardsByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	userRepoMock.EXPECT().
		GetBoardsByID(testUser.ID).
		Return(nil, nil, nil)

	admins, members, err := uUsecase.GetBoardsByID(testUser.ID)
	assert.Nil(t, admins)
	assert.Nil(t, members)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	sid := "test_sid"
	sessionRepoMock.EXPECT().
		Delete("test_sid").
		Return(nil)
	userRepoMock.EXPECT().
		Delete(testUser.ID).
		Return(nil)

	deleteErr := uUsecase.Delete(testUser.ID, sid)

	assert.NoError(t, deleteErr)

}
