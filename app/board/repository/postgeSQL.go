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

func (boardStore *BoardStore) Create(board *models.Board) error {
	return boardStore.DB.Create(board).Error
}

func (boardStore *BoardStore) Get(bid uint) *models.Board {
	brd := new(models.Board)
	brd.ID = bid
	err := boardStore.DB.Model(brd).Related(&brd.Admins, "Admins").Error
	if err != nil {
		return nil
	}
	for _, admins := range brd.Admins {
		admins.Password = ""
	}
	err = boardStore.DB.Model(brd).Related(&brd.Members, "Members").Error
	if err != nil {
		return nil
	}
	for _, members := range brd.Members {
		members.Password = ""
	}
	return brd
}

func (boardStore *BoardStore) GetAll(usr *models.User) ([]models.Board, []models.Board, error) {
	var adminsBoard []models.Board
	err := boardStore.DB.Model(usr).Related(&adminsBoard, "Admin").Error
	if err != nil {
		return nil, nil, err
	}
	var membersBoard []models.Board
	err = boardStore.DB.Model(usr).Related(&membersBoard, "Member").Error
	if err != nil {
		return nil, nil, err
	}
	usr.Password = ""
	return adminsBoard, membersBoard, nil
}
