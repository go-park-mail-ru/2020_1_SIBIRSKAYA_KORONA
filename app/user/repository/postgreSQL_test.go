package repository_test

import (
	"flag"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/repository"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"os"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var test_opts struct {
	configPath string
}

func TestMain(m *testing.M) {
	flag.StringVar(&test_opts.configPath, "test-c", "", "path to configuration file")
	flag.StringVar(&test_opts.configPath, "test-config", "", "path to configuration file")
	flag.Parse()

	viper.SetConfigFile(test_opts.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	logger.InitLogger()
	os.Exit(m.Run())
}

func SetupDB() (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("cant create mock: %s", err)
	}
	DB, erro := gorm.Open("postgres", db)
	// DB.AutoMigrate(&models.Board{})
	if erro != nil {
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
		Email:    "email",
		Avatar:   "avatar",
		Password: []byte("password"),
	}

	repo := repository.CreateRepository(db)

	// good
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "users"`).WithArgs(
		usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Email, usr.Avatar, sqlmock.AnyArg()).WillReturnRows(
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
		WithArgs(usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Email, usr.Avatar, sqlmock.AnyArg()).
		WillReturnError(errors.ErrConflict)

	if err := repo.Create(&usr); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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
		[]string{"id", "name", "surname", "nickname", "email", "avatar"}).AddRow(
		usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Email, usr.Avatar))
	getUsr, err := repo.GetByID(usr.ID)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if usr.ID != getUsr.ID || usr.Name != getUsr.Name || usr.Surname != getUsr.Surname ||
		usr.Nickname != getUsr.Nickname || usr.Email != getUsr.Email || usr.Avatar != getUsr.Avatar {
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
		[]string{"id", "name", "surname", "nickname", "email", "avatar"}).AddRow(
		usr.ID, usr.Name, usr.Surname, usr.Nickname, usr.Email, usr.Avatar))
	getUsr, err := repo.GetByNickname(usr.Nickname)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
		return
	}
	if usr.ID != getUsr.ID || usr.Name != getUsr.Name || usr.Surname != getUsr.Surname ||
		usr.Nickname != getUsr.Nickname || usr.Email != getUsr.Email || usr.Avatar != getUsr.Avatar {
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
