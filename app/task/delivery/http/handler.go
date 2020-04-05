package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/task"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/message"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	useCase task.UseCase
}

func CreateHandler(router *echo.Echo, useCase task.UseCase, mw *middleware.GoMiddleware) {
	handler := &TaskHandler{useCase: useCase}
	router.POST("boards/:bid/columns/:cid/tasks", handler.Create,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard)
	router.GET("boards/:bid/columns/:cid/tasks/:tid", handler.Get,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard)
	router.PUT("boards/:bid/columns/:cid/tasks/:tid", handler.Update,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	router.DELETE("boards/:bid/columns/:cid/tasks/:tid", handler.Delete,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
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
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "task", Data: *tsk})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, body)
}

func (taskHandler *TaskHandler) Get(ctx echo.Context) error {
	var cid, tid uint
	_, err := fmt.Sscan(ctx.Param("cid"), &cid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	_, err = fmt.Sscan(ctx.Param("tid"), &tid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	tsk, useErr := taskHandler.useCase.Get(cid, tid)
	if useErr != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(useErr), message.ResponseError{Message: useErr.Error()})
	}
	body, err := message.GetBody(message.Pair{Name: "task", Data: *tsk})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, body)
}

func (taskHandler *TaskHandler) Update(ctx echo.Context) error {
	tsk := models.CreateTask(ctx)
	if tsk == nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if _, err := fmt.Sscan(ctx.Param("tid"), &tsk.ID); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := taskHandler.useCase.Update(*tsk)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	return ctx.NoContent(http.StatusOK)
}

func (taskHandler *TaskHandler) Delete(ctx echo.Context) error {
	var tid uint
	if _, err := fmt.Sscan(ctx.Param("tid"), &tid); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err := taskHandler.useCase.Delete(tid)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(errors.ResolveErrorToCode(err), message.ResponseError{Message: err.Error()})
	}
	return ctx.NoContent(http.StatusOK)
}
