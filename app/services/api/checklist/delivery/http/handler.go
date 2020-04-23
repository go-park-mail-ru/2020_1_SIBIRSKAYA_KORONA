package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/labstack/echo/v4"
)

type ChecklistHandler struct {
	useCase checklist.UseCase
}

func CreateHandler(router *echo.Echo, useCase checklist.UseCase, mw *middleware.GoMiddleware) {
	handler := &ChecklistHandler{
		useCase: useCase,
	}
	router.GET("/boards/:bid/columns/:cid/tasks/:tid/checklists", handler.Get,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	router.POST("/boards/:bid/columns/:cid/tasks/:tid/checklists", handler.Create,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	router.PUT("/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid", handler.Update,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckChecklistInTask)
	router.DELETE("/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid", handler.Delete,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckChecklistInTask)
}

func (checklistHandler *ChecklistHandler) Create(ctx echo.Context) error {
	var chlist models.Checklist
	body := ctx.Get("body").([]byte)
	err := chlist.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	chlist.Tid = ctx.Get("tid").(uint)
	err = checklistHandler.useCase.Create(&chlist)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := chlist.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (checklistHandler *ChecklistHandler) Get(ctx echo.Context) error {
	tid := ctx.Get("tid").(uint)
	chlists, err := checklistHandler.useCase.Get(tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := chlists.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (checklistHandler *ChecklistHandler) Update(ctx echo.Context) error {
	return ctx.NoContent(errors.ResolveErrorToCode(errors.ErrDbBadOperation))
}

func (checklistHandler *ChecklistHandler) Delete(ctx echo.Context) error {
	var clid uint
	_, err := fmt.Sscan(ctx.Param("clid"), &clid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if checklistHandler.useCase.Delete(clid) != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}
