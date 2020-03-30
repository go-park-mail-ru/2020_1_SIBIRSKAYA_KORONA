package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"
	"net/http"

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
	// CORS
	router.OPTIONS("/boards", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})
	router.OPTIONS("/boards/:bid", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})

	router.POST("/boards", handler.Create)
	router.GET("/boards", handler.GetAll)

	router.GET("/boards/:bid", handler.Get)
	router.PUT("/boards/:bid", handler.Update)
	router.DELETE("/boards/:bid", handler.Delete)

	router.GET("/boards/:bid/members", handler.throwError)
	router.POST("/boards/:bid/members", handler.throwError)

	router.GET("/boards/:bid/tasks", handler.throwError)
	router.POST("/boards/:bid/tasks", handler.throwError)

	router.GET("/boards/:bid/labels", handler.throwError)
	router.POST("/boards/:bid/labels", handler.throwError)

	// TODO(Alexandr | Timofei): move to label handler
	//router.GET("/boards/:bid/labels/:lid", handler.throwError)
	//router.PUT("/boards/:bid/labels/:lid", handler.throwError)
	//router.DELETE("/boards/:bid/labels/:lid", handler.throwError)
}

// TODO(Alexandr): remove after debug
func (boardHandler *BoardHandler) throwError(ctx echo.Context) error {
	panic("handler not implemented")
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

func (boardHandler *BoardHandler) GetAll(ctx echo.Context) error {
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	bAdmin, bMember, err := boardHandler.useCase.GetAll(cookie.Value)
	// TODO: Антон
	if err != nil {
		return ctx.NoContent(http.StatusNotFound)
	}
	body, err := message.GetBody(message.Pair{Name: "admin", Data: bAdmin}, message.Pair{Name: "member", Data: bMember})
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}

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

func (boardHandler *BoardHandler) Update(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (boardHandler *BoardHandler) Delete(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
