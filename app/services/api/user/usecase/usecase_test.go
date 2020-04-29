package usecase_test

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/mocks"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"os"
	"testing"

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
		Return("cookie_value", nil)

	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)
	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	sid, err := uUsecase.Create(&testUser, 3)
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

	// good
	ID := testUser.ID
	userRepoMock.EXPECT().
		GetByID(ID).
		Return(&testUser, nil)

	user, err := uUsecase.GetByID(ID)
	assert.NoError(t, err)
	assert.Equal(t, &testUser, user)

	// error
	ID++
	userRepoMock.EXPECT().
		GetByID(ID).
		Return(nil, errors.ErrUserNotFound)
	_, err = uUsecase.GetByID(ID)
	assert.Equal(t, err, errors.ErrUserNotFound)
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

	// good
	nickname := testUser.Nickname
	userRepoMock.EXPECT().
		GetByNickname(nickname).
		Return(&testUser, nil)

	user, err := uUsecase.GetByNickname(nickname)
	assert.NoError(t, err)
	assert.Equal(t, &testUser, user)

	// error
	nickname += "a"
	userRepoMock.EXPECT().
		GetByNickname(nickname).
		Return(nil, errors.ErrUserNotFound)
	_, err = uUsecase.GetByNickname(nickname)
	assert.Equal(t, err, errors.ErrUserNotFound)
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

	oldPass := []byte("oldPass")
	userRepoMock.EXPECT().
		Update([]byte("oldPass"), testUser, nil).
		Return(nil)
	updateErr := uUsecase.Update(oldPass, testUser, nil)

	assert.NoError(t, updateErr)

	// error not found
	testUser.ID++
	userRepoMock.EXPECT().
		Update(nil, testUser, nil).
		Return(errors.ErrUserNotFound)
	err = uUsecase.Update(nil, testUser, nil)
	assert.Equal(t, err, errors.ErrUserNotFound)

	// error wrong pass
	testUser.ID--
	oldPass = []byte("dfvdfvdfv")
	userRepoMock.EXPECT().
		Update(oldPass, testUser, nil).
		Return(errors.ErrWrongPassword)
	err = uUsecase.Update(oldPass, testUser, nil)
	assert.Equal(t, err, errors.ErrWrongPassword)
}

//GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error)
func TestGetUsersByNicknamePart(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUsers models.Users
	err := faker.FakeData(&testUsers)
	assert.NoError(t, err)

	nicknamePart := "aaa"
	var limit uint = 3
	// good
	userRepoMock.EXPECT().GetUsersByNicknamePart(nicknamePart, limit).Return(testUsers, nil)
	usrs, err := uUsecase.GetUsersByNicknamePart(nicknamePart, limit)
	assert.NoError(t, err)
	assert.Equal(t, len(testUsers), len(usrs))
	for idx := range testUsers {
		assert.Equal(t, usrs[idx], testUsers[idx])
	}

	// error not found
	nicknamePart += "a"
	limit--
	userRepoMock.EXPECT().GetUsersByNicknamePart(nicknamePart, limit).Return(nil, errors.ErrUserNotFound)
	_, err = uUsecase.GetUsersByNicknamePart(nicknamePart, limit)
	assert.Equal(t, err, errors.ErrUserNotFound)
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

	// good
	sid := "test_sid"
	sessionRepoMock.EXPECT().Delete(sid).Return(nil)
	userRepoMock.EXPECT().Delete(testUser.ID).Return(nil)
	err = uUsecase.Delete(testUser.ID, sid)
	assert.NoError(t, err)

	// err NoCookie
	sid = "test_sid1"
	sessionRepoMock.EXPECT().Delete(sid).Return(errors.ErrNoCookie)
	err = uUsecase.Delete(testUser.ID, sid)
	assert.Equal(t, err, errors.ErrNoCookie)

	// err not found
	sid = "test_sid"
	testUser.ID++
	sessionRepoMock.EXPECT().Delete(sid).Return(nil)
	userRepoMock.EXPECT().Delete(testUser.ID).Return(errors.ErrUserNotFound)
	err = uUsecase.Delete(testUser.ID, sid)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

/*func TestGetBoardsByID(t *testing.T) {
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
}*/
/*func TestCheckPassword(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, sessionRepoMock := createRepoMocks(ctrl)
	uUsecase := userUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var id uint = 1
	pass := []byte("12345678")
	// good
	userRepoMock.EXPECT().CheckPassword(id, pass).Return(true)
	ok := uUsecase.//CheckPassword(id, pass)
	assert.NoError(t, err)
}*/
