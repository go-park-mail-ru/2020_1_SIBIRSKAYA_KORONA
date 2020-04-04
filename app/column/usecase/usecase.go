package usecase

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
)

type ColumnUseCase struct {
	columnRepo column.Repository
}

func CreateUseCase(columnRepo_ column.Repository) column.UseCase {
	return &ColumnUseCase{columnRepo: columnRepo_}
}

func (columnUseCase *ColumnUseCase) Create(column *models.Column) error {
	return columnUseCase.columnRepo.Create(column)
}
