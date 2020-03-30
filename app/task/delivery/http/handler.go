package http

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	useCase task.UseCase
}

func CreateHandler(router *echo.Echo, useCase task.UseCase) {
	handler := &TaskHandler{useCase: useCase}

	router.GET("boards/:bid/columns/:cid/tasks", handler.throwError)
	router.POST("boards/:bid/columns/:cid/tasks", handler.throwError)

	router.GET("boards/:bid/tasks/:tid", handler.throwError)
	router.PUT("boards/:bid/tasks/:tid", handler.throwError)
	router.DELETE("boards/:bid/tasks/:tid", handler.throwError)

	router.GET("/boards/:bid/tasks/:tid/labels", handler.throwError)
	router.POST("/boards/:bid/tasks/:tid/labels", handler.throwError)
	router.DELETE("/boards/:bid/tasks/:tid/labels/:lid", handler.throwError)

	router.GET("/boards/:bid/tasks/:tid/members", handler.throwError)
	router.POST("/boards/:bid/tasks/:tid/members", handler.throwError)
	router.DELETE("/boards/:bid/tasks/:tid/members/:uid", handler.throwError)
}

// TODO(Alexandr): remove after debug
func (columnHandler *TaskHandler) throwError(ctx echo.Context) error {
	panic("handler not implemented")
}
