package repository_test

import (
	"log"
	"os"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/repository"
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

	var tsk models.Task
	err := faker.FakeData(&tsk)
	assert.NoError(t, err)

	repo := repository.CreateRepository(db)

	// good
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "tasks"`).WithArgs(
		tsk.ID, tsk.Name, tsk.About, tsk.Level, tsk.Deadline, tsk.Pos, tsk.Cid).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(tsk.ID))
	mock.ExpectCommit()

	if err := repo.Create(&tsk); err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "tasks"`).WithArgs(
		tsk.ID, tsk.Name, tsk.About, tsk.Level, tsk.Deadline, tsk.Pos, tsk.Cid).WillReturnError(errors.ErrConflict)

	if err := repo.Create(&tsk); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	t.Skip()
	t.Parallel()

	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()

	var tsk models.Task
	err := faker.FakeData(&tsk)
	assert.NoError(t, err)

	repo := repository.CreateRepository(db)

	mock.ExpectQuery(`SELECT (\*) FROM (.*)"tasks" WHERE (.*)"tasks"."id" (.*) LIMIT 1`).WithArgs(
		tsk.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "name", "about", "level", "deadline", "pos", "cid"}).AddRow(
		tsk.ID, tsk.Name, tsk.About, tsk.Level, tsk.Deadline, tsk.Pos, tsk.Cid))

	// newUsr := usr
	// newUsr.Name = "name1"
	// newUsr.Surname = "name1"
	// newUsr.Nickname = "name1"
	// newUsr.Email = "name1"
	// newUsr.Password = nil

	var newTsk models.Task
	err = faker.FakeData(&newTsk)
	assert.NoError(t, err)
	newTsk.ID = tsk.ID

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.*)"tasks" SET (.*) WHERE (.*)"tasks"`).WithArgs(
		tsk.ID, tsk.Name, tsk.About, tsk.Level, tsk.Deadline, tsk.Pos, tsk.Cid).WillReturnResult(
		sqlmock.NewResult(int64(newTsk.ID), 1))
	mock.ExpectCommit()

	if err := repo.Update(newTsk); err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// // error
	// usr.ID++
	// mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
	// 	usr.ID).WillReturnError(errors.ErrUserNotFound)
	// if err := repo.Update(nil, usr); err == nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// 	return
	// }
	// if err := mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// 	return
	// }
}
