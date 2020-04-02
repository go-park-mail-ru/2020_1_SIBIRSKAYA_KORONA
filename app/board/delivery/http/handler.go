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

type ResponseError struct {
	Message string `json:"message"`
}

func CreateHandler(router *echo.Echo, useCase board.UseCase, mw *middleware.GoMiddleware) {
	handler := &BoardHandler{
		useCase: useCase,
	}
	router.OPTIONS("/boards", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})
	router.POST("/boards", handler.Create, mw.CheckCookieExist)
	router.GET("/boards/:bid/members", handler.Get, mw.CheckCookieExist)
	router.GET("/boards/:bid/tasks", handler.Get, mw.CheckCookieExist)
	router.GET("/boards", handler.GetAll, mw.CheckCookieExist)
	router.PUT("/boards", handler.Update)
}

func (boardHandler *BoardHandler) Create(ctx echo.Context) error {
	cookie := ctx.Get("sid").(string)

	brd := models.CreateBoard(ctx)
	if brd == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if boardHandler.useCase.Create(cookie, brd) != nil {
		return ctx.NoContent(http.StatusInternalServerError) // TODO: пока хз
	}

	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	cookie := ctx.Get("sid").(string)

	var bid uint
	_, err := fmt.Sscan(ctx.Param("bid"), &bid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	brd, useErr := boardHandler.useCase.Get(cookie, bid)
	if useErr != nil {
		return ctx.JSON(errors.ResolveErrorToCode(useErr), ResponseError{Message: useErr.Error()})
	}

	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) GetAll(ctx echo.Context) error {
	cookie := ctx.Get("sid").(string)

	bAdmin, bMember, useErr := boardHandler.useCase.GetAll(cookie)

	if useErr != nil {
		return ctx.JSON(errors.ResolveErrorToCode(useErr), ResponseError{Message: useErr.Error()})
	}

	body, err := message.GetBody(message.Pair{Name: "admin", Data: bAdmin}, message.Pair{Name: "member", Data: bMember})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Update(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
