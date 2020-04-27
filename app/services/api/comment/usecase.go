package comment

import "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

//go:generate mockgen -source=usecase.go -package=mocks -destination=./mocks/comment_usecase_mock.go
type UseCase interface {
	CreateComment(comment *models.Comment) error
	GetComments(tid uint, uid uint) (models.Comments, error)
	GetByID(tid uint, comid uint) (*models.Comment, error)
	Delete(comid uint) error
}
