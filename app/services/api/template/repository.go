package template

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type Repository interface {
	Create(tmpl *models.Template) error
}
