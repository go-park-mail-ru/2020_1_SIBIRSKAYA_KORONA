package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	commentMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/mocks"
	commentUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/usecase"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createRepoMocks(controller *gomock.Controller) *commentMocks.MockRepository {
	comRepoMock := commentMocks.NewMockRepository(controller)
	return comRepoMock
}

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func TestCreateComment(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepoMock := createRepoMocks(ctrl)
	comUsecase := commentUseCase.CreateUseCase(commentRepoMock)

	var testComment models.Comment
	err := faker.FakeData(&testComment)
	assert.NoError(t, err)

	commentRepoMock.EXPECT().
		CreateComment(&testComment).
		Return(nil)

	err = comUsecase.CreateComment(&testComment)
	assert.NoError(t, err)

	commentRepoMock.EXPECT().
		CreateComment(&testComment).
		Return(errors.ErrDbBadOperation)

	err = comUsecase.CreateComment(&testComment)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())
}

func TestGetComments(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepoMock := createRepoMocks(ctrl)
	comUsecase := commentUseCase.CreateUseCase(commentRepoMock)

	var comments models.Comments
	var testComment models.Comment
	err := faker.FakeData(&testComment)
	assert.NoError(t, err)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)

	var testUser models.User
	err = faker.FakeData(&testUser)
	assert.NoError(t, err)

	testComment.Uid = testUser.ID
	comments = append(comments, testComment)

	commentRepoMock.EXPECT().
		GetComments(testTask.ID).
		Return(comments, nil)

	cmts, err := comUsecase.GetComments(testTask.ID, testUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, cmts, comments)

	commentRepoMock.EXPECT().
		GetComments(testTask.ID).
		Return(nil, errors.ErrTaskNotFound)

	cmts, err = comUsecase.GetComments(testTask.ID, testUser.ID)
	assert.Nil(t, cmts)
	assert.EqualError(t, err, errors.ErrTaskNotFound.Error())
}

func TestGetByID(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepoMock := createRepoMocks(ctrl)
	comUsecase := commentUseCase.CreateUseCase(commentRepoMock)

	var testComment models.Comment
	err := faker.FakeData(&testComment)
	assert.NoError(t, err)

	var testTask models.Task
	err = faker.FakeData(&testTask)
	assert.NoError(t, err)

	commentRepoMock.EXPECT().
		GetByID(testComment.ID).
		Return(&testComment, nil)

	testComment.Tid = testTask.ID + 1

	cmt, err := comUsecase.GetByID(testTask.ID, testComment.ID)
	assert.Nil(t, cmt)
	assert.EqualError(t, err, errors.ErrNoPermission.Error())

	testComment.Tid = testTask.ID

	commentRepoMock.EXPECT().
		GetByID(testComment.ID).
		Return(&testComment, nil)

	cmt, err = comUsecase.GetByID(testTask.ID, testComment.ID)
	assert.NoError(t, err)
	assert.Equal(t, cmt, &testComment)

	commentRepoMock.EXPECT().
		GetByID(testComment.ID).
		Return(nil, errors.ErrCommentNotFound)

	cmt, err = comUsecase.GetByID(testTask.ID, testComment.ID)
	assert.Nil(t, cmt)
	assert.EqualError(t, err, errors.ErrCommentNotFound.Error())
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	commentRepoMock := createRepoMocks(ctrl)
	comUsecase := commentUseCase.CreateUseCase(commentRepoMock)

	var testComment models.Comment
	err := faker.FakeData(&testComment)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	commentRepoMock.EXPECT().
		Delete(testComment.ID).
		Return(nil)

	err = comUsecase.Delete(testComment.ID)
	assert.NoError(t, err)

	commentRepoMock.EXPECT().
		Delete(testComment.ID).
		Return(errors.ErrDbBadOperation)

	err = comUsecase.Delete(testComment.ID)
	assert.EqualError(t, err, errors.ErrDbBadOperation.Error())
}
