package http

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	useCase task.UseCase
}

func CreateHandler(router *echo.Echo, useCase task.UseCase) {
	handler:= &TaskHandler{useCase:useCase}

	router.GET("boards/:bid/tasks/:tid", handler.throwError)
	router.PUT("boards/:bid/tasks/:tid", handler.throwError)
	router.DELETE("boards/:bid/tasks/:tid", handler.throwError)
}

// TODO(Alexandr): remove after debug
func (columnHandler *TaskHandler) throwError(ctx echo.Context) error {
	panic("handler not implemented")
}