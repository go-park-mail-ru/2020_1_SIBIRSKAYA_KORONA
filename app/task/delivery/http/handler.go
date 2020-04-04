package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TaskHandler struct {
	useCase task.UseCase
}

func CreateHandler(router *echo.Echo, useCase task.UseCase, mw *middleware.GoMiddleware) {
	handler := &TaskHandler{useCase: useCase}
	router.POST("boards/:bid/columns/:cid/tasks", handler.Create,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard)

	/*router.GET("boards/:bid/tasks/:tid", handler.throwError)
	router.PUT("boards/:bid/tasks/:tid", handler.throwError)
	router.DELETE("boards/:bid/tasks/:tid", handler.throwError)

	router.GET("/boards/:bid/tasks/:tid/labels", handler.throwError)
	router.POST("/boards/:bid/tasks/:tid/labels", handler.throwError)
	router.DELETE("/boards/:bid/tasks/:tid/labels/:lid", handler.throwError)

	router.GET("/boards/:bid/tasks/:tid/members", handler.throwError)
	router.POST("/boards/:bid/tasks/:tid/members", handler.throwError)
	router.DELETE("/boards/:bid/tasks/:tid/members/:uid", handler.throwError)*/
}

func (taskHandler *TaskHandler) Create(ctx echo.Context) error {
	tsk := models.CreateTask(ctx)
	if tsk == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if _, err := fmt.Sscan(ctx.Param("cid"), &tsk.Cid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := taskHandler.useCase.Create(tsk)
	if err != nil {
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "task", Data: *tsk})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, body)
}
