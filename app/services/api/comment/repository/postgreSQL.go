package repository

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type CommentStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) comment.Repository {
	return &CommentStore{DB: db}
}

func (commentStore *CommentStore) CreateComment(cmt *models.Comment) error {
	var user models.User
	err := commentStore.DB.Select("nickname, avatar").
		Where("id = ?", cmt.Uid).
		Find(&user).Error

	if err != nil {
		logger.Error(err)
		return errors.ErrUserNotFound
	}

	cmt.Avatar = user.Avatar
	cmt.Nickname = user.Nickname
	err = commentStore.DB.Create(cmt).Error
	if err != nil {
		logger.Error(err)
		//return errors.ErrConflict
		return errors.ErrDbBadOperation
	}

	return nil
}

func (commentStore *CommentStore) GetComments(tid uint) (models.Comments, error) {
	var cmts models.Comments
	err := commentStore.DB.Model(&models.Task{ID: tid}).Related(&cmts, "tid").Error
	// err := columnStore.DB.Model(&models.Column{ID: cid}).Preload("Members").Related(&tsks, "cid").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrDbBadOperation
	}

	var user models.User
	// TODO: попробовать через preload
	for id := range cmts {
		err := commentStore.DB.Select("nickname, avatar").
			Where("id = ?", cmts[id].Uid).
			Find(&user).Error
		if err != nil {
			logger.Error(err)
			return nil, errors.ErrUserNotFound
		}
		cmts[id].Avatar = user.Avatar
		cmts[id].Nickname = user.Nickname
	}

	return cmts, nil
}

func (commentStore *CommentStore) GetByID(comid uint) (*models.Comment, error) {
	comment := new(models.Comment)
	if err := commentStore.DB.First(comment, comid).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrCommentNotFound
	}

	return comment, nil
}

func (commentStore *CommentStore) Delete(comid uint) error {
	err := commentStore.DB.Delete(&models.Comment{ID: comid}).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}
