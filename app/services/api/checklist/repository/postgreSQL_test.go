package repository_test

import (
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist/repository"
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

	var chlist models.Checklist
	err := faker.FakeData(&chlist)
	assert.NoError(t, err)

	repo := repository.CreateRepository(db)

	// good
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "checklists"`).WithArgs(
		chlist.ID, chlist.Name, chlist.Tid).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(chlist.ID))
	mock.ExpectCommit()

	if err := repo.Create(&chlist); err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "checklists"`).WithArgs(
		chlist.ID, chlist.Name, chlist.Tid).WillReturnError(errors.ErrConflict)

	if err := repo.Create(&chlist); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	t.Skip()
	t.Parallel()

	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	repo := repository.CreateRepository(db)

	// good
	var checklist models.Checklist
	err := faker.FakeData(&checklist)
	assert.NoError(t, err)

	var item models.Item
	err = faker.FakeData(&item)
	assert.NoError(t, err)

	item.Clid = checklist.Tid

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "checklists" WHERE (id = $1)`)).WithArgs(checklist.ID).WillReturnResult(
		sqlmock.NewResult(int64(checklist.ID), 1))
	mock.ExpectCommit()

	query := regexp.QuoteMeta(`SELECT * FROM "items" WHERE ("clid" = $1)`)
	mock.ExpectQuery(query).WithArgs(
		checklist.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "text", "is_done", "clid"}).AddRow(
		item.ID, item.Text, item.IsDone, item.Clid))

	if err := repo.Delete(checklist.ID); err != nil {
		t.Fatalf("unexpected error: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// // error
	// mock.ExpectBegin()
	// mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "checklists" WHERE (id = $1)`)).WithArgs(checklist.ID).
	// 	WillReturnError(errors.ErrUserNotFound)
	// mock.ExpectCommit()

	// if err := repo.Delete(checklist.ID); err == nil {
	// 	t.Errorf("expected error, got nil")
	// }
	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)

	// }

}
