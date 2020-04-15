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
		return errors.ErrConflict
	}
	return nil
}

func (boardStore *BoardStore) Get(bid uint) (*models.Board, error) {
	brd := new(models.Board)
	err := boardStore.DB.First(brd, bid).Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrBoardNotFound
	}
	err = boardStore.DB.Model(brd).Related(&brd.Admins, "Admins").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	err = boardStore.DB.Model(brd).Related(&brd.Members, "Members").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}
	return brd, nil
}

func (boardStore *BoardStore) GetColumnsByID(bid uint) ([]models.Column, error) {
	var cols []models.Column
	err := boardStore.DB.Model(&models.Board{ID: bid}).Related(&cols, "bid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrBoardNotFound
	}
	return cols, nil
}

func (boardStore *BoardStore) Update(newBoard *models.Board) error {
	oldBoard := new(models.Board)
	err := boardStore.DB.First(oldBoard, newBoard.ID).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrBoardNotFound
	}
	oldBoard.Name = newBoard.Name
	err = boardStore.DB.Save(oldBoard).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}
	return nil
}

func (boardStore *BoardStore) Delete(bid uint) error {
	err := boardStore.DB.Delete(&models.Column{ID: bid}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrBoardNotFound
	}
	return nil
}

func (boardStore *BoardStore) InviteMember(bid uint, member *models.User) error {
	brd := new(models.Board)
	err := boardStore.DB.First(brd, bid).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrBoardNotFound
	}

	err = boardStore.DB.Model(&brd).Association("Members").Append(member).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrBoardNotFound
	}

	return nil
}

func (boardStore *BoardStore) DeleteMember(bid uint, member *models.User) error {
	brd := new(models.Board)
	err := boardStore.DB.First(brd, bid).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrBoardNotFound
	}

	err = boardStore.DB.Model(&brd).Association("Members").Delete(member).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrBoardNotFound
	}

	return nil
}
