package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/labstack/echo/v4"
)

type BoardHandler struct {
	useCase board.UseCase
}

func CreateHandler(router *echo.Echo, useCase board.UseCase, mw *middleware.GoMiddleware) {
	handler := &BoardHandler{
		useCase: useCase,
	}
	router.POST("/boards", handler.Create, mw.CheckAuth)
	router.GET("/boards/:bid", handler.Get, mw.CheckAuth)
	router.GET("/boards/:bid/columns", handler.GetColumns, mw.CheckAuth, mw.CheckBoardMemberPermission)
	router.PUT("/boards/:bid", handler.Update, mw.CheckAuth, mw.CheckBoardAdminPermission)
	router.DELETE("/boards/:bid", handler.Delete, mw.CheckAuth, mw.CheckBoardAdminPermission) // TODO: что если есть другие админы

	//router.GET("/boards/:bid/labels", handler.throwError)
	//router.POST("/boards/:bid/labels", handler.throwError)
	//router.GET("/boards/:bid/labels/:lid", handler.throwError)
	//router.PUT("/boards/:bid/labels/:lid", handler.throwError)
	//router.DELETE("/boards/:bid/labels/:lid", handler.throwError)
}

func (boardHandler *BoardHandler) Create(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	brd := models.CreateBoard(ctx)
	if brd == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := boardHandler.useCase.Create(uid, brd)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	var bid uint
	_, err := fmt.Sscan(ctx.Param("bid"), &bid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	brd, err := boardHandler.useCase.Get(uid, bid, false)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) GetColumns(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	cols, err := boardHandler.useCase.GetColumnsByID(bid)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "columns", Data: cols})
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Update(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (boardHandler *BoardHandler) Delete(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
