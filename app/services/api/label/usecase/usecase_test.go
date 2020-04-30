package usecase_test

import (
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	labelMocks "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/mocks"
	labelUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/usecase"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createRepoMocks(controller *gomock.Controller) *labelMocks.MockRepository {
	itemRepoMock := labelMocks.NewMockRepository(controller)
	return itemRepoMock
}

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)

	labelRepoMock.EXPECT().
		Create(&testLabel).
		Return(nil)

	err = lUsecase.Create(&testLabel)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	// t.Skip()
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	labelRepoMock := createRepoMocks(ctrl)
	lUsecase := labelUseCase.CreateUseCase(labelRepoMock)

	var testLabel models.Label
	err := faker.FakeData(&testLabel)
	assert.NoError(t, err)
	//t.Logf("%+v", testBoard)

	labelRepoMock.EXPECT().
		Delete(testLabel.ID).
		Return(nil)

	err = lUsecase.Delete(testLabel.ID)
	assert.NoError(t, err)
}
