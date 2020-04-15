package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/labstack/echo/v4"
)

type ColumnHandler struct {
	useCase column.UseCase
}

func CreateHandler(router *echo.Echo, useCase column.UseCase, mw *middleware.GoMiddleware) {
	handler := &ColumnHandler{useCase: useCase}
	router.POST("/boards/:bid/columns", handler.Create, mw.CheckAuth, mw.CheckBoardAdminPermission)
	router.GET("/boards/:bid/columns/:cid", handler.Get, mw.CheckAuth, mw.CheckBoardMemberPermission)
	router.GET("/boards/:bid/columns/:cid/tasks", handler.GetTasks,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard)
	router.PUT("/boards/:bid/columns/:cid", handler.Update,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard)
	router.DELETE("/boards/:bid/columns/:cid", handler.Delete,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard)
}

func (columnHandler *ColumnHandler) Create(ctx echo.Context) error {
	col := models.CreateColumn(ctx)
	if col == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	col.Bid = ctx.Get("bid").(uint)
	err := columnHandler.useCase.Create(col)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "column", Data: *col})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, body)
}

func (columnHandler *ColumnHandler) Get(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	var cid uint
	_, err := fmt.Sscan(ctx.Param("cid"), &cid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	col, err := columnHandler.useCase.Get(bid, cid)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "column", Data: *col})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (columnHandler *ColumnHandler) GetTasks(ctx echo.Context) error {
	cid := ctx.Get("cid").(uint)
	tsks, err := columnHandler.useCase.GetTasksByID(cid)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "tasks", Data: tsks})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (columnHandler *ColumnHandler) Update(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

// TODO: проверить удаление колонки, если в ней есть таски
func (columnHandler *ColumnHandler) Delete(ctx echo.Context) error {
	cid := ctx.Get("cid").(uint)
	err := columnHandler.useCase.Delete(cid)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	return ctx.NoContent(http.StatusOK)
}
