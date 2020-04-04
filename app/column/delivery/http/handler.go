package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ColumnHandler struct {
	useCase column.UseCase
}

func CreateHandler(router *echo.Echo, useCase column.UseCase, mw *middleware.GoMiddleware) {
	handler := &ColumnHandler{useCase: useCase}
	router.POST("/boards/:bid/columns", handler.Create, mw.AuthByCookie, mw.CheckBoardAdminPermission)
	router.GET("/boards/:bid/columns/:cid", handler.throwError)
	router.PUT("/boards/:bid/columns/:cid", handler.throwError)
	router.DELETE("/boards/:bid/columns/:cid", handler.throwError)
}

func (columnHandler *ColumnHandler) Create(ctx echo.Context) error {
	col := models.CreateColumn(ctx)
	if col == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if _, err := fmt.Sscan(ctx.Param("bid"), &col.Bid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := columnHandler.useCase.Create(col)
	if err != nil {
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, col)
}

func (columnHandler *ColumnHandler) throwError(ctx echo.Context) error {
	panic("handler not implemented")
}
