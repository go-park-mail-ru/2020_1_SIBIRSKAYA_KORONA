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
	// CORS
	router.OPTIONS("/boards", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})
	router.OPTIONS("/boards/:bid", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})

    // TODO: ОФОРМИТЬ В ГРУППУ
	router.POST("/boards", handler.Create, mw.CheckCookieExist)
	router.GET("/boards", handler.GetAll, mw.CheckCookieExist)

	router.GET("/boards/:bid", handler.Get, mw.CheckCookieExist)
	router.PUT("/boards/:bid", handler.Update, mw.CheckCookieExist)
	router.DELETE("/boards/:bid", handler.Delete, mw.CheckCookieExist)

	router.GET("/boards/:bid/members", handler.throwError, mw.CheckCookieExist)
	router.POST("/boards/:bid/members", handler.throwError, mw.CheckCookieExist)
	router.DELETE("/boards/:bid/members/:uid", handler.throwError, mw.CheckCookieExist)

	router.GET("/boards/:bid/admins", handler.throwError, mw.CheckCookieExist)
	router.POST("/boards/:bid/admins", handler.throwError, mw.CheckCookieExist)
	router.DELETE("/boards/:bid/admins/:uid", handler.throwError, mw.CheckCookieExist)



	// TODO(Alexandr | Timofey): move to label handler

	//router.GET("/boards/:bid/labels", handler.throwError)
	//router.POST("/boards/:bid/labels", handler.throwError)
	//router.GET("/boards/:bid/labels/:lid", handler.throwError)
	//router.PUT("/boards/:bid/labels/:lid", handler.throwError)
	//router.DELETE("/boards/:bid/labels/:lid", handler.throwError)
}

// TODO(Alexandr): remove after debug
func (boardHandler *BoardHandler) throwError(ctx echo.Context) error {
	panic("handler not implemented")
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

func (boardHandler *BoardHandler) Delete(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
