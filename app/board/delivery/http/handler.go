package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
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
	//admin := router.Group("/boards")
	//admin.Use()
	router.POST("/boards", handler.Create, mw.CheckAuth)
	router.GET("/boards/:bid", handler.Get, mw.CheckAuth)
	router.GET("/boards/:bid/columns", handler.GetColumns, mw.CheckAuth, mw.CheckBoardMemberPermission)

	//router.PUT("/boards/:bid", handler.Update, mw.CheckBoardAdminPermission, mw.CheckAuth)
	//router.DELETE("/boards/:bid", handler.Delete, mw.CheckBoardAdminPermission, mw.CheckAuth)

	//router.GET("/boards/:bid/labels", handler.throwError)
	//router.POST("/boards/:bid/labels", handler.throwError)
	//router.GET("/boards/:bid/labels/:lid", handler.throwError)
	//router.PUT("/boards/:bid/labels/:lid", handler.throwError)
	//router.DELETE("/boards/:bid/labels/:lid", handler.throwError)
}

func (boardHandler *BoardHandler) Create(ctx echo.Context) error {
	userID := ctx.Get("userID").(uint)
	brd := models.CreateBoard(ctx)
	if brd == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if boardHandler.useCase.Create(userID, brd) != nil {
		return ctx.NoContent(http.StatusInternalServerError) // TODO: пока хз
	}
	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	userID := ctx.Get("userID").(uint)
	var bid uint
	_, err := fmt.Sscan(ctx.Param("bid"), &bid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	brd, useErr := boardHandler.useCase.Get(userID, bid, false)
	if useErr != nil {
		return ctx.JSON(errors.ResolveErrorToCode(useErr), message.ResponseError{Message: useErr.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) GetColumns(ctx echo.Context) error {
	var bid uint
	_, err := fmt.Sscan(ctx.Param("bid"), &bid)
	cols, useErr := boardHandler.useCase.GetColumnsByID(bid)
	if useErr != nil {
		return ctx.JSON(errors.ResolveErrorToCode(useErr), message.ResponseError{Message: useErr.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "columns", Data: cols})
	if err != nil {
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
