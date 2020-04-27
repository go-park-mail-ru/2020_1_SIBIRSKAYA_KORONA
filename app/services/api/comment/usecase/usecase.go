package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type CommentUseCase struct {
	commentRepo comment.Repository
}

func CreateUseCase(commentRepo_ comment.Repository) comment.UseCase {
	return &CommentUseCase{commentRepo: commentRepo_}
}

func (commentUseCase *CommentUseCase) CreateComment(cmt *models.Comment) error {
	err := commentUseCase.commentRepo.CreateComment(cmt)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (commentUseCase *CommentUseCase) GetComments(tid uint, readerID uint) (models.Comments, error) {
	cmts, err := commentUseCase.commentRepo.GetComments(tid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	for cmtID := range cmts {
		if cmts[cmtID].Uid == readerID {
			cmts[cmtID].ReaderIsAuthor = true
		}
	}
	return cmts, nil
}

func (commentUseCase *CommentUseCase) GetByID(tid uint, comid uint) (*models.Comment, error) {
	comment, err := commentUseCase.commentRepo.GetByID(comid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if comment.Tid != tid {
		return nil, errors.ErrNoPermission
	}
	return comment, nil
}

func (commentUseCase *CommentUseCase) Delete(comid uint) error {
	err := commentUseCase.commentRepo.Delete(comid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
