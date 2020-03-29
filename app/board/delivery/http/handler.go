package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/labstack/echo/v4"
)

type BoardHandler struct {
	useCase board.UseCase
}

func CreateHandler(router *echo.Echo, useCase board.UseCase) {
	handler := &BoardHandler{
		useCase: useCase,
	}
	router.OPTIONS("/boards", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})
	router.POST("/boards", handler.Create)
	router.GET("/boards/:bid/members", handler.Get)
	router.GET("/boards/:bid/tasks", handler.Get)
	router.GET("/boards", handler.GetAll)
	router.PUT("/boards", handler.Update)
}

func (boardHandler *BoardHandler) Create(ctx echo.Context) error {
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//

	brd := models.CreateBoard(ctx)
	if brd == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if boardHandler.useCase.Create(cookie.Value, brd) != nil {
		return ctx.NoContent(http.StatusInternalServerError) // TODO: пока хз
	}
	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//
	var bid uint
	_, err = fmt.Sscan(ctx.Param("bid"), &bid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	brd := boardHandler.useCase.Get(cookie.Value, bid)
	if brd == nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	body, err := message.GetBody(message.Pair{Name: "board", Data: *brd})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) GetAll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}

	bAdmin, bMember, useErr := boardHandler.useCase.GetAll(cookie.Value)

	if useErr.Err != nil {
		return ctx.JSON(useErr.Code, useErr.Err.Error())
	}

	body, err := message.GetBody(message.Pair{Name: "admin", Data: bAdmin}, message.Pair{Name: "member", Data: bMember})
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Update(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
