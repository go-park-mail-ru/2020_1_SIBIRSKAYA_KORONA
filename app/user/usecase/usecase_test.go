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

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := userMocks.NewMockRepository(ctrl)
	sessionRepoMock := sessionMocks.NewMockRepository(ctrl)

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
