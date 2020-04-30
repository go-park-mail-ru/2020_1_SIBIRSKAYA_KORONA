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

func TestGetBoardsByUser(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var id uint = 1

	// good
	boardRepoMock.EXPECT().GetBoardsByUser(id).Return(nil, nil, nil)
	_, _, err := bUsecase.GetBoardsByUser(id)
	assert.Nil(t, err)

	// error
	id++
	boardRepoMock.EXPECT().GetBoardsByUser(id).Return(nil, nil, errors.ErrBoardNotFound)
	_, _, err = bUsecase.GetBoardsByUser(id)
	assert.Equal(t, err, errors.ErrBoardNotFound)
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

	// good
	boardRepoMock.EXPECT().Delete(testBoard.ID).Return(nil)
	err = bUsecase.Delete(testBoard.ID)
	assert.NoError(t, err)

	// error
	testBoard.ID++
	boardRepoMock.EXPECT().Delete(testBoard.ID).Return(errors.ErrUserNotFound)
	err = bUsecase.Delete(testBoard.ID)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestGetColumnsByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testBoard models.Board
	err := faker.FakeData(&testBoard)
	assert.NoError(t, err)

	// good
	boardRepoMock.EXPECT().GetColumnsByID(testBoard.ID).Return(nil, nil)
	column, err := bUsecase.GetColumnsByID(testBoard.ID)
	assert.Nil(t, column)
	assert.NoError(t, err)

	// error
	boardRepoMock.EXPECT().GetColumnsByID(testBoard.ID).Return(nil, errors.ErrColNotFound)
	column, err = bUsecase.GetColumnsByID(testBoard.ID)
	assert.Nil(t, column)
	assert.Equal(t, err, errors.ErrColNotFound)
}

func TestGetLabelsByID(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var bid uint = 1

	// good
	boardRepoMock.EXPECT().GetLabelsByID(bid).Return(nil, nil)
	column, err := bUsecase.GetLabelsByID(bid)
	assert.Nil(t, column)
	assert.NoError(t, err)

	// error
	boardRepoMock.EXPECT().GetLabelsByID(bid).Return(nil, errors.ErrLabelNotFound)
	column, err = bUsecase.GetLabelsByID(bid)
	assert.Nil(t, column)
	assert.Equal(t, err, errors.ErrLabelNotFound)
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

	// good
	boardRepoMock.EXPECT().Update(&testBoard).Return(nil)
	err = bUsecase.Update(&testBoard)
	assert.NoError(t, err)

	// error
	testBoard.ID++
	boardRepoMock.EXPECT().Update(&testBoard).Return(errors.ErrBoardNotFound)
	err = bUsecase.Update(&testBoard)
	assert.Equal(t, err, errors.ErrBoardNotFound)
}

func TestGetUsersForInvite(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var bid uint = 0
	var limit uint = 0
	nicknamePart := "aaa"

	// good
	boardRepoMock.EXPECT().GetUsersForInvite(bid, nicknamePart, limit).Return(nil, nil)
	_, err := bUsecase.GetUsersForInvite(bid, nicknamePart, limit)
	assert.NoError(t, err)

	// error
	bid++
	boardRepoMock.EXPECT().GetUsersForInvite(bid, nicknamePart, limit).Return(nil, errors.ErrUserNotFound)
	_, err = bUsecase.GetUsersForInvite(bid, nicknamePart, limit)
	assert.Equal(t, err, errors.ErrUserNotFound)
}

func TestInviteMember(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	var bid uint

	// good
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	boardRepoMock.EXPECT().InviteMember(bid, &testUser).Return(nil)
	err = bUsecase.InviteMember(bid, testUser.ID)
	assert.NoError(t, err)

	// error user not found
	testUser.ID++
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
	err = bUsecase.InviteMember(bid, testUser.ID)
	assert.Equal(t, err, errors.ErrUserNotFound)

	// board not found
	testUser.ID--
	bid++
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	boardRepoMock.EXPECT().InviteMember(bid, &testUser).Return(errors.ErrBoardNotFound)
	err = bUsecase.InviteMember(bid, testUser.ID)
	assert.Equal(t, err, errors.ErrBoardNotFound)
}

func TestDeleteMember(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock, boardRepoMock := createRepoMocks(ctrl)
	bUsecase := boardUseCase.CreateUseCase(userRepoMock, boardRepoMock)

	var testUser models.User
	err := faker.FakeData(&testUser)
	assert.NoError(t, err)
	var bid uint

	// good
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	boardRepoMock.EXPECT().DeleteMember(bid, &testUser).Return(nil)
	err = bUsecase.DeleteMember(bid, testUser.ID)
	assert.NoError(t, err)

	// error user not found
	testUser.ID++
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(nil, errors.ErrUserNotFound)
	err = bUsecase.DeleteMember(bid, testUser.ID)
	assert.Equal(t, err, errors.ErrUserNotFound)

	// board not found
	testUser.ID--
	bid++
	userRepoMock.EXPECT().GetByID(testUser.ID).Return(&testUser, nil)
	boardRepoMock.EXPECT().DeleteMember(bid, &testUser).Return(errors.ErrBoardNotFound)
	err = bUsecase.DeleteMember(bid, testUser.ID)
	assert.Equal(t, err, errors.ErrBoardNotFound)
}
