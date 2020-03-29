package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/cstmerr"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
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
	err := boardStore.DB.Model(brd).Related(&brd.Admins, "Admins").Related(&brd.Members, "Members").Error
	if err != nil {
		return nil
	}
	for _, member := range append(brd.Admins, brd.Members...) {
		member.Password = ""
	}
	return brd
}

func (boardStore *BoardStore) GetAll(usr *models.User) ([]models.Board, []models.Board, *cstmerr.RepoError) {
	var adminsBoard, membersBoard []models.Board
	err := boardStore.DB.Model(usr).Related(&adminsBoard, "Admin").Related(&membersBoard, "Member").Error
	if err != nil {
		return nil, nil, &cstmerr.RepoError{Err: errors.Wrap(err, models.ErrDbBadOperation.Error())}
	}
	return adminsBoard, membersBoard, &cstmerr.RepoError{Err: nil}
}
