package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

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

func (boardStore *BoardStore) Get(bid uint) (*models.Board, error) {
	brd := new(models.Board)
	brd.ID = bid
	err := boardStore.DB.Model(brd).Related(&brd.Admins, "Admins").Related(&brd.Members, "Members").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}

	for _, member := range append(brd.Admins, brd.Members...) {
		member.Password = ""
	}
	return brd, nil
}

func (boardStore *BoardStore) GetAll(usr *models.User) ([]models.Board, []models.Board, error) {
	var adminsBoard, membersBoard []models.Board

	err := boardStore.DB.Model(usr).Related(&adminsBoard, "Admin").Related(&membersBoard, "Member").Error
	if err != nil {
		logger.Error(err)
		return nil, nil, errors.ErrDbBadOperation
	}

	return adminsBoard, membersBoard, nil
}
