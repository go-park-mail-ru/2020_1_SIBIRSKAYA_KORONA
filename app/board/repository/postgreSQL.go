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
	err := boardStore.DB.Create(board).Error
	if err != nil {
	    logger.Error(err)
		return errors.ErrDbBadOperation
	return nil
}


func (boardStore *BoardStore) Get(bid uint) (*models.Board, error) {
	brd := new(models.Board)
	brd.ID = bid
	err := boardStore.DB.First(brd).Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	return brd, nil
}

func (boardStore *BoardStore) GetAll(usr *models.User) ([]models.Board, []models.Board, error) {
	var adminsBoards, membersBoards []models.Board
	err := boardStore.DB.Model(usr).
		Related(&adminsBoards, "Admin").
		Related(&membersBoards, "Member").Error

	if err != nil {
	    logger.Error(err)
		return nil, nil, errors.ErrDbBadOperation
	}
	return adminsBoards, membersBoards, nil
}

func (boardStore *BoardStore) Update(newBoard *models.Board) error {
	oldBoard := new(models.Board)
	err := boardStore.DB.First(oldBoard, newBoard.ID).Error
	if err != nil {
	    logger.Error(err)
		return errors.ErrDbBadOperation
    }
	oldBoard.Name = newBoard.Name

	err = boardStore.DB.Save(oldBoard).Error
	if err != nil {
	    logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (boardStore *BoardStore) Delete(bid uint) error {
	err := boardStore.DB.Delete(&models.Column{ID: bid}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	} 
	return nil
}
