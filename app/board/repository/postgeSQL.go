package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type BoardStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) board.Repository {
	return &BoardStore{DB: db}
}

func (userStore *BoardStore) Create(board *models.Board) error {
	return userStore.DB.Create(board).Error
}

func (userStore *BoardStore) GetAll(usr *models.User) ([]models.Board, []models.Board, error) {
	var adminsBoard []models.Board
	err := userStore.DB.Model(usr).Related(&adminsBoard, "Admin").Error
	if err != nil {
		return nil, nil, err
	}
	var membersBoard []models.Board
	err = userStore.DB.Model(usr).Related(&membersBoard, "Member").Error
	if err != nil {
		return nil, nil, err
	}
	return adminsBoard, membersBoard, nil
}
