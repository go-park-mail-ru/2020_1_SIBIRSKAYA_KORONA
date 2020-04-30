package usecase_test

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	boardMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board/mocks"
	boardUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board/usecase"
	userMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/mocks"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	os.Exit(m.Run())
}

func createRepoMocks(controller *gomock.Controller) (*userMocks.MockRepository, *boardMocks.MockRepository) {
	userRepoMock := userMocks.NewMockRepository(controller)
	boardRepoMock := boardMocks.NewMockRepository(controller)
	return userRepoMock, boardRepoMock
}

func TestCreate(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)

	// good
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	boardRepoMock.EXPECT().Create(testUser.ID, &testBoard).Return(nil)
	err = bUsecase.Create(testUser.ID, &testBoard)
	assert.NoError(t, err)

	// error user not found
	testUser.ID++
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
	err = bUsecase.Create(testUser.ID, &testBoard)
	assert.Equal(t, err, errors.ErrUserNotFound)

	// error conflict
	testUser.ID--
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	boardRepoMock.EXPECT().Create(testUser.ID, &testBoard).Return(errors.ErrConflict)
	err = bUsecase.Create(testUser.ID, &testBoard)
	assert.Equal(t, err, errors.ErrConflict)
}

func TestGet(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)

	var testBoard models.Board
	err = faker.FakeData(&testBoard)
	assert.NoError(t, err)
	/*err = faker.FakeData(testBoard.Members)
	assert.NoError(t, err)
	for idx := range testBoard.Admins {
		if testBoard.Admins[idx].ID == testUser.ID {
			testBoard.Admins[idx].ID++
		}
	}
	log.Println(testBoard)*/
	// no permission admin
	boardRepoMock.EXPECT().Get(testBoard.ID).Return(&testBoard, nil)
	board, err := bUsecase.Get(testUser.ID, testBoard.ID, true)
	assert.Nil(t, board)
	assert.Equal(t, err, errors.ErrNoPermission)
	// good
	/* testBoard.Admins[0].ID = testUser.ID
	boardRepoMock.EXPECT().Get(testBoard.ID).Return(&testBoard, nil)
	board, err = bUsecase.Get(testUser.ID, testBoard.ID, true)
	assert.NoError(t, err)
	assert.Equal(t, board, testBoard)*/
}

func TestDelete(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testBoard models.Board
	err := faker.FakeData(&testBoard)
	assert.NoError(t, err)

	boardRepoMock.EXPECT().Delete(testBoard.ID).Return(nil)

	err = bUsecase.Delete(testBoard.ID)
	assert.NoError(t, err)
}

/*

func TestGetColumnsByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testBoard models.Board
	err := faker.FakeData(&testBoard)
	assert.NoError(t, err)

	boardRepoMock.EXPECT().
		GetColumnsByID(testBoard.ID).
		Return(nil, nil)

	column, err := bUsecase.GetColumnsByID(testBoard.ID)
	assert.Nil(t, column)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testBoard models.Board
	err := faker.FakeData(&testBoard)
	assert.NoError(t, err)

	boardRepoMock.EXPECT().
		Update(&testBoard).
		Return(nil)

	err = bUsecase.Update(&testBoard)
	assert.NoError(t, err)
}
*/
