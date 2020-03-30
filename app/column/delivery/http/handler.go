package http

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column"
	"github.com/labstack/echo/v4"
)

type ColumnHandler struct {
	useCase column.UseCase
}

func CreateHandler(router *echo.Echo, useCase column.UseCase) {
	handler:= &ColumnHandler{useCase:useCase}

	router.GET("boards/:bid/columns/:cid", handler.throwError)
	router.PUT("boards/:bid/columns/:cid", handler.throwError)
	router.DELETE("boards/:bid/columns/:cid", handler.throwError)
}

// TODO(Alexandr): remove after debug
func (columnHandler *ColumnHandler) throwError(ctx echo.Context) error {
	panic("handler not implemented")
}