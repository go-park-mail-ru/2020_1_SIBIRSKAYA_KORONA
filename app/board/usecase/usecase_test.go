package usecase_test

import (
	"flag"
	"os"
	"testing"

	"github.com/bxcodec/faker"
	boardMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board/mocks"
	boardUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/mocks"
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

func createRepoMocks(controller *gomock.Controller) (*userMocks.MockRepository, *boardMocks.MockRepository) {
	userRepoMock := userMocks.NewMockRepository(controller)
	boardRepoMock := boardMocks.NewMockRepository(controller)
	return userRepoMock, boardRepoMock
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	userRepoMock.EXPECT().
		GetByID(testUser.ID).
		Return(&testUser, nil)

	boardRepoMock.EXPECT().
		Create(&testBoard).
		Return(nil)

	err = bUsecase.Create(testUser.ID, &testBoard)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	boardRepoMock.EXPECT().
		Get(testBoard.ID).
		Return(&testBoard, nil)

	board, err := bUsecase.Get(testUser.ID, testBoard.ID, true)

	assert.Nil(t, board)
	assert.Equal(t, err, errors.ErrNoPermission)
}

func TestGetColumnsByID(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testBoard models.Board
	err := faker.FakeData(&testBoard)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	boardRepoMock.EXPECT().
		GetColumnsByID(testBoard.ID).
		Return(nil, nil)

	column, err := bUsecase.GetColumnsByID(testBoard.ID)
	assert.Nil(t, column)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testBoard models.Board
	err := faker.FakeData(&testBoard)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	boardRepoMock.EXPECT().
		Update(gomock.Any).
		Return(nil)

	err = bUsecase.Update(testBoard)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {

}
