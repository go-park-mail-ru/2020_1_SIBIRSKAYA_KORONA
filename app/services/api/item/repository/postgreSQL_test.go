package repository_test

import (
	"log"
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/repository"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestMain(m *testing.M) {
	logger.InitLoggerByConfig(logger.LoggerConfig{Logfile: "stdout", Loglevel: zapcore.DebugLevel})
	os.Exit(m.Run())
}

func SetupDB() (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("cant create mock: %s", err)
	}
	DB, err := gorm.Open("postgres", db)
	if err != nil {
		log.Fatalf("Got an unexpected error: %s", err)
	}
	return mock, DB
}

func TestCreate(t *testing.T) {
	// t.Skip()
	t.Parallel()

	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()

	var itm models.Item
	err := faker.FakeData(&itm)
	assert.NoError(t, err)

	repo := repository.CreateRepository(db)

	// good
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "items"`).WithArgs(
		itm.ID, itm.Text, itm.IsDone, itm.Clid).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(itm.ID))
	mock.ExpectCommit()

	if err := repo.Create(&itm); err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "items"`).WithArgs(
		itm.ID, itm.Text, itm.IsDone, itm.Clid).WillReturnError(errors.ErrConflict)

	if err := repo.Create(&itm); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
