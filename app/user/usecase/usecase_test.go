package usecase_test

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/mocks"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/mocks"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/usecase"
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

func createRepoMocks(controller *gomock.Controller) (*userMocks.MockRepository, *sessionMocks.MockRepository) {
	userRepoMock := userMocks.NewMockRepository(controller)
	sessionRepoMock := sessionMocks.NewMockRepository(controller)
	return userRepoMock, sessionRepoMock
}

func TestCreate(t *testing.T) {
	// t.Skip()
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

	oldPass := "oldPass"

	// передать nil в качестве *multipart.FileHeader приемлемо, это самый частый кейс
	// работу с картинкой будем гонять на тестах юзерского репозитория
	userRepoMock.EXPECT().
		Update("oldPass", &testUser, nil).
		Return(nil)

	updateErr := uUsecase.Update(oldPass, &testUser, nil)

	assert.NoError(t, updateErr)

	// второго вызова к юзерской репе произойти не должно
	updateErr = uUsecase.Update("no", nil, nil)

	assert.Equal(t, updateErr, errors.ErrInternal)

}
