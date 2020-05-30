package repository_test

import (
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/repository"
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

	var attach models.AttachedFile
	err := faker.FakeData(&attach)
	assert.NoError(t, err)

	repo := repository.CreateRepository(db)

	// good
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "attached_files"`).WithArgs(
		attach.ID, attach.URL, attach.Name, attach.FileKey, attach.Tid).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(attach.ID))
	mock.ExpectCommit()

	if err := repo.Create(&attach); err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "attached_files"`).WithArgs(
		attach.ID, attach.URL, attach.Name, attach.FileKey, attach.Tid).WillReturnError(errors.ErrConflict)

	if err := repo.Create(&attach); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {
	//t.Skip()
	t.Parallel()

	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	repo := repository.CreateRepository(db)

	// good
	var attach models.AttachedFile
	err := faker.FakeData(&attach)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "attached_files" WHERE (id = $1)`)).WithArgs(attach.ID).WillReturnResult(
		sqlmock.NewResult(int64(attach.ID), 1))

	mock.ExpectCommit()
	if err := repo.Delete(attach.ID); err != nil {
		t.Fatalf("unexpected error: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "attached_files" WHERE (id = $1)`)).WithArgs(attach.ID).
		WillReturnError(errors.ErrUserNotFound)
	if err := repo.Delete(attach.ID); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestGetByID(t *testing.T) {
	//t.Skip()
	t.Parallel()

	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()

	var attach models.AttachedFile
	err := faker.FakeData(&attach)
	assert.NoError(t, err)

	repo := repository.CreateRepository(db)

	// good

	query := regexp.QuoteMeta(`SELECT * FROM "attached_files" WHERE (id = $1) ORDER BY "attached_files"."id" ASC LIMIT 1`)

	mock.ExpectQuery(query).WithArgs(
		attach.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "url", "name", "filekey", "tid"}).AddRow(
		attach.ID, attach.URL, attach.Name, attach.FileKey, attach.Tid))

	if _, err := repo.GetByID(attach.ID); err != nil {
		t.Errorf("unexpected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	query = regexp.QuoteMeta(`SELECT * FROM "attached_files" WHERE (id = $1) ORDER BY "attached_files"."id" ASC LIMIT 1`)
	// error
	attach.ID++
	mock.ExpectQuery(query).WithArgs(
		attach.ID).WillReturnError(errors.ErrFileNotFound)
	if _, err := repo.GetByID(attach.ID); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
