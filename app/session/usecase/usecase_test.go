package usecase_test

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/mocks"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/usecase"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/mocks"
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
	sUsecase := sessionUseCase.CreateUseCase(sessionRepoMock, userRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testUser)

	sessionExpires := time.Now().AddDate(1, 0, 0)

	userRepoMock.EXPECT().
		GetByNickname(testUser.Nickname).
		Return(&testUser, nil)

	sessionRepoMock.EXPECT().
		Create(gomock.Any()).
		Return("cookie_value", nil)

	sid, err := sUsecase.Create(&testUser, sessionExpires)

	assert.NoError(t, err)
	assert.Equal(t, sid, "cookie_value")
}

func TestGet(t *testing.T) {
	// t.Skip()
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
	// t.Skip()
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
