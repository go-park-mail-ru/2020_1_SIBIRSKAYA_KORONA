package template

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type Usecase interface {
	Create(uid uint, tmpl *models.Template) (*models.Board, error)
}
