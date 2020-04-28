package comment

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/comment_repo_mock.go
type Repository interface {
	CreateComment(comment *models.Comment) error
	GetComments(tid uint) (models.Comments, error)
	GetByID(comid uint) (*models.Comment, error)
	Delete(comid uint) error
}
