package repository_test

import (
	"log"
	"os"
	"regexp"
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

	var firstMember models.User
	err = faker.FakeData(&firstMember)
	assert.NoError(t, err)

	var secondMember models.User
	err = faker.FakeData(&secondMember)
	assert.NoError(t, err)

	var thirdMember models.User
	err = faker.FakeData(&thirdMember)
	assert.NoError(t, err)

	var firstLabel models.Label
	err = faker.FakeData(&firstLabel)
	assert.NoError(t, err)

	repo := repository.CreateRepository(db)

	mock.ExpectQuery(`SELECT (\*) FROM (.*)"tasks" WHERE (.*)"tasks"."id" (.*) LIMIT 1`).WithArgs(tsk.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "about", "level", "deadline", "pos", "cid"}).AddRow(
			tsk.ID, tsk.Name, tsk.About, tsk.Level, tsk.Deadline, tsk.Pos, tsk.Cid))

	query := regexp.QuoteMeta(`SELECT "users".* FROM "users" INNER JOIN "task_members" ON "task_members"."user_id" = "users"."id" WHERE ("task_members"."task_id" IN ($1))`)
	mock.ExpectQuery(query).WithArgs(tsk.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "name", "surname", "nickname", "avatar", "email"}).
		AddRow(firstMember.ID, firstMember.Name, firstMember.Surname, firstMember.Nickname, firstMember.Avatar, firstMember.Email).
		AddRow(secondMember.ID, secondMember.Name, secondMember.Surname, secondMember.Nickname, secondMember.Avatar, secondMember.Email).
		AddRow(thirdMember.ID, thirdMember.Name, thirdMember.Surname, thirdMember.Nickname, thirdMember.Avatar, thirdMember.Email))

	query = regexp.QuoteMeta(`SELECT "labels".* FROM "labels" INNER JOIN "task_labels" ON "task_labels"."label_id" = "labels"."id" WHERE ("task_labels"."task_id" IN ($1)) ORDER BY "id"`)
	mock.ExpectQuery(query).WithArgs(tsk.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "name", "color", "bid"}).
		AddRow(firstLabel.ID, firstLabel.Name, firstLabel.Color, firstLabel.Bid))

	var newTsk models.Task
	err = faker.FakeData(&newTsk)
	assert.NoError(t, err)
	newTsk.ID = tsk.ID

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.*)"tasks" SET (.*) WHERE (.*)"tasks"`).WithArgs(
		newTsk.Name, newTsk.About, newTsk.Level, newTsk.Deadline, newTsk.Pos, newTsk.Cid, newTsk.ID).WillReturnResult(
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

func TestDelete(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	repo := repository.CreateRepository(db)

	// good
	var tsk models.Task
	err := faker.FakeData(&tsk)
	assert.NoError(t, err)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE (id = $1)`)).WithArgs(tsk.ID).WillReturnResult(
		sqlmock.NewResult(int64(tsk.ID), 1))

	mock.ExpectCommit()
	if err := repo.Delete(tsk.ID); err != nil {
		t.Fatalf("unexpected error: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE (id = $1)`)).WithArgs(tsk.ID).
		WillReturnError(errors.ErrUserNotFound)
	if err := repo.Delete(tsk.ID); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
