package repository_test

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/repository"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	pass "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/password"

	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
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
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	usr := models.User{
		ID:       1,
		Name:     "name",
		Surname:  "surname",
		Nickname: "nickname",
		Avatar:   "avatar",
		Email:    "email",
		Password: []byte("password"),
	}

	repo := repository.CreateRepository(db)

	// good
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "users"`).WithArgs(
		usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Avatar, usr.Email, sqlmock.AnyArg()).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(usr.ID))
	mock.ExpectCommit()

	if err := repo.Create(&usr); err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// error
	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO (.*) "users"`).
		WithArgs(usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Avatar, usr.Email, sqlmock.AnyArg()).
		WillReturnError(errors.ErrConflict)

	if err := repo.Create(&usr); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func UserEqual(l, r models.User) bool {
	return l.ID == r.ID && l.Name == r.Name && l.Surname == r.Surname &&
		l.Nickname == r.Nickname && l.Email == r.Email && l.Avatar == r.Avatar
}

func TestGetByID(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	usr := models.User{
		ID:       104,
		Name:     "name",
		Surname:  "surname",
		Nickname: "nickname",
		Email:    "email",
		Avatar:   "avatar",
	}
	repo := repository.CreateRepository(db)

	// good
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		usr.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "name", "surname", "nickname", "avatar", "email"}).AddRow(
		usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Avatar, usr.Email))
	getUsr, err := repo.GetByID(usr.ID)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if !UserEqual(usr, *getUsr) {
		t.Errorf("results not match, want %v, have %v", usr, *getUsr)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	usr.ID++
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		usr.ID).WillReturnError(errors.ErrUserNotFound)
	if _, err := repo.GetByID(usr.ID); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByNickName(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	usr := models.User{
		ID:       14,
		Name:     "name",
		Surname:  "surname",
		Nickname: "lovelove",
		Email:    "email",
		Avatar:   "avatar",
	}
	repo := repository.CreateRepository(db)

	// good
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE \(nickname = (.*)\) ORDER BY "users"."id" ASC LIMIT 1`).WithArgs(
		usr.Nickname).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "name", "surname", "nickname", "avatar", "email"}).AddRow(
		usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Avatar, usr.Email))
	getUsr, err := repo.GetByNickname(usr.Nickname)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if !UserEqual(usr, *getUsr) {
		t.Errorf("results not match, want %v, have %v", usr, *getUsr)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	usr.Nickname = usr.Nickname + "aaa"
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE \(nickname = (.*)\) ORDER BY "users"."id" ASC LIMIT 1`).WithArgs(
		usr.Nickname).WillReturnError(errors.ErrUserNotFound)
	if _, err := repo.GetByNickname(usr.Nickname); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	usr := models.User{
		ID:       904,
		Name:     "name",
		Surname:  "surname",
		Nickname: "lovelovelove",
		Email:    "email",
		Avatar:   "avatar",
		Password: []byte("aaa"),
	}

	repo := repository.CreateRepository(db)
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		usr.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "name", "surname", "nickname", "avatar", "email", "password"}).AddRow(
		usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Avatar, usr.Email, usr.Password))
	newUsr := usr
	newUsr.Name = "name1"
	newUsr.Surname = "name1"
	newUsr.Nickname = "name1"
	newUsr.Email = "name1"
	newUsr.Password = nil

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE (.*)"users" SET (.*) WHERE (.*)"users"`).WithArgs(
		newUsr.Name, newUsr.Surname, newUsr.Nickname, newUsr.Avatar, newUsr.Email, usr.Password, newUsr.ID).WillReturnResult(
		sqlmock.NewResult(int64(newUsr.ID), 1))
	mock.ExpectCommit()

	if err := repo.Update(nil, newUsr); err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	usr.ID++
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		usr.ID).WillReturnError(errors.ErrUserNotFound)
	if err := repo.Update(nil, usr); err == nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestCheckPassword(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	usr := models.User{
		ID:       904,
		Password: []byte("lovelove"),
	}
	hashPass := pass.HashPasswordGenSalt(usr.Password)
	repo := repository.CreateRepository(db)

	// good
	mock.ExpectQuery(`SELECT password FROM "users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		usr.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(usr.ID, hashPass))
	if ok := repo.CheckPassword(usr.ID, usr.Password); !ok {
		t.Fatal("unexpected error", ok)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	mock.ExpectQuery(`SELECT password FROM "users" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		usr.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow(usr.ID, hashPass))
	usr.Password = []byte("azazaz")
	if ok := repo.CheckPassword(usr.ID, usr.Password); ok {
		t.Fatal("unexpected error", ok)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestDelete(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	repo := repository.CreateRepository(db)

	// good
	var id uint = 10
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE (id = $1)`)).WithArgs(id).WillReturnResult(
		sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
	if err := repo.Delete(id); err != nil {
		t.Fatalf("unexpected error: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE (id = $1)`)).WithArgs(id).
		WillReturnError(errors.ErrUserNotFound)
	if err := repo.Delete(id); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetUsersByNicknamePart(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	usr1 := models.User{
		ID:       14,
		Name:     "name",
		Surname:  "surname",
		Nickname: "lovelove",
		Email:    "email",
		Avatar:   "avatar",
		Password: []byte("aaa"),
	}
	usr2 := models.User{
		ID:       15,
		Name:     "name",
		Surname:  "surname",
		Nickname: "lovelo",
		Email:    "email",
		Avatar:   "avatar",
		Password: []byte("aaa"),
	}
	usr3 := models.User{
		ID:       16,
		Name:     "name",
		Surname:  "surname",
		Nickname: "lovel",
		Email:    "email",
		Avatar:   "avatar",
		Password: []byte("aaa"),
	}
	repo := repository.CreateRepository(db)

	// good
	part := "love"
	var limit uint = 3
	query := regexp.QuoteMeta(`SELECT * FROM "users" WHERE (nickname LIKE $1) LIMIT ` + fmt.Sprintf("%d", limit))
	mock.ExpectQuery(query).WithArgs(part + "%").WillReturnRows(sqlmock.NewRows(
		[]string{"id", "name", "surname", "nickname", "avatar", "email"}).
		AddRow(usr1.ID, usr1.Name, usr1.Surname, usr1.Nickname, usr1.Avatar, usr1.Email).
		AddRow(usr2.ID, usr2.Name, usr2.Surname, usr2.Nickname, usr2.Avatar, usr2.Email).
		AddRow(usr3.ID, usr3.Name, usr3.Surname, usr3.Nickname, usr3.Avatar, usr3.Email))
	usrs, err := repo.GetUsersByNicknamePart(part, limit)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if !UserEqual(usr1, usrs[0]) || !UserEqual(usr2, usrs[1]) || !UserEqual(usr3, usrs[2]) {
		t.Errorf("wrong answer")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error
	part = "aaaaaaaaaaaaaaaaaaaa"
	query = regexp.QuoteMeta(`SELECT * FROM "users" WHERE (nickname LIKE $1) LIMIT ` + fmt.Sprintf("%d", limit))
	mock.ExpectQuery(query).WithArgs(part + "%").WillReturnError(errors.ErrUserNotFound)
	usrs, err = repo.GetUsersByNicknamePart(part, limit)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if len(usrs) != 0 || err != nil {
		t.Errorf("wrong answer")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}
