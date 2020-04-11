package repository_test

import (
	"flag"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board/repository"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	// "github.com/stretchr/testify/require"

	//"database/sql"
	//"fmt"
	// "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"log"
	"os"
	"testing"

	// "github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	// "reflect"
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
	if erro != nil {
		log.Fatalf("Got an unexpected error: %s", err)

	}
	return mock, DB
}

func TestCreate(t *testing.T) {
	mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()
	brd := models.Board{ID:1, Name:"name1"}
	repo := repository.CreateRepository(db)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO (.*) "boards"`).WithArgs(
		brd.ID, brd.Name).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(brd.ID))
	mock.ExpectCommit()
	if err := repo.Create(&brd); err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	brd.Name = ""
	mock.ExpectBegin()
	mock.
		ExpectQuery(`INSERT INTO (.*) "boards"`).
		WithArgs(brd.ID, brd.Name).WillReturnError(errors.ErrConflict)

	if err := repo.Create(&brd); err == nil {
		t.Errorf("expected error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGet(t *testing.T) {
	/*mock, db := SetupDB()
	defer db.Close()
	defer mock.ExpectClose()

	brd := models.Board{ID:1, Name:"name1"}
	mock.ExpectQuery(`SELECT (\*) FROM (.*)"boards" WHERE (.*)"users"."id" (.*) LIMIT 1`).WithArgs(
		brd.ID).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "username", "password", "email", "image"}).AddRow(
		brd.ID, brd.Name, brd.Password, brd.Email, testUser.Image))

	if usr, err := ud.GetById(testUser.Id); err != nil {
		t.Fatalf("unexpected error %s", err)
	} else {
		assert.Equal(t, testUser, *usr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}*/
}
