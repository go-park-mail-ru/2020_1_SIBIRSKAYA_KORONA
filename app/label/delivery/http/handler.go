package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/labstack/echo/v4"
)

type BoardHandler struct {
	useCase board.UseCase
}

func CreateHandler(router *echo.Echo, useCase board.UseCase, mw *middleware.GoMiddleware) {
	handler := &BoardHandler{
		useCase: useCase,
	}
	//router.GET("/boards/:bid/labels", handler.Get)
	//router.POST("/boards/:bid/labels", handler.throwError)
	router.GET("/boards/:bid/labels/:lid", handler.Get)
	//router.PUT("/boards/:bid/labels/:lid", handler.throwError)
	//router.DELETE("/boards/:bid/labels/:lid", handler.throwError)
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
