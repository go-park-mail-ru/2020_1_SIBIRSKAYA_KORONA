package http

import (
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
	router.OPTIONS("/boards", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})
	router.POST("/boards", handler.Create)
	router.GET("/boards", handler.Get)
	router.GET("/boards/:board", handler.GetAll)
	router.PUT("/boards", handler.Update)
}

/*
	/boards (GET)  -  все доски. дл каждой доски: админ, имя доски, участники, BID
	/boards (POST) -  с фронта: имя доски. с бэка статус + bid
	/boards/{bid} (GET) -  конкретная доска: админ, имя доски, uid - участников, tid - задач
*/

func (boardHandler *BoardHandler) Create(ctx echo.Context) error {
	// в миддлвар
	cookie, err := ctx.Cookie("session_id")
	if err != nil {
		return ctx.NoContent(http.StatusForbidden)
	}
	//

	board := models.CreateBoard(ctx)
	if board == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if boardHandler.useCase.Create(cookie.Value, board) != nil {
		return ctx.NoContent(http.StatusInternalServerError) // TODO: пока хз
	}
	body, err := message.GetBody(message.Pair{Name: "board", Data: *board})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (boardHandler *BoardHandler) GetAll(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func (boardHandler *BoardHandler) Update(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}