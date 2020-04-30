package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	commentMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/mocks"
	commentUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/usecase"

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
}
