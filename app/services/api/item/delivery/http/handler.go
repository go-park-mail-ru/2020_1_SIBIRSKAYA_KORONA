package http

import (
	"net/http"

	"fmt"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	useCase item.UseCase
}

func CreateHandler(router *echo.Echo, useCase item.UseCase, mw *middleware.GoMiddleware) {
	handler := &ItemHandler{
		useCase: useCase,
	}

	router.POST("/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid/items", handler.Create,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckChecklistInTask)
	router.PUT("/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid/items/:itid", handler.Update,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckChecklistInTask, mw.CheckItemInChecklist)
	router.DELETE("/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid/items/:itid", handler.Delete,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckChecklistInTask, mw.CheckItemInChecklist)
}

func (itemHandler *ItemHandler) Create(ctx echo.Context) error {
	var clid uint
	_, err := fmt.Sscan(ctx.Param("clid"), &clid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var item models.Item
	body := ctx.Get("body").([]byte)
	err = item.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	item.Clid = clid

	err = itemHandler.useCase.Create(&item)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := item.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (itemHandler *ItemHandler) Update(ctx echo.Context) error {
	clid := ctx.Get("clid").(uint)

	itid := ctx.Get("itid").(uint)

	var item models.Item
	body := ctx.Get("body").([]byte)
	err := item.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	item.Clid = clid
	item.ID = itid

	err = itemHandler.useCase.Update(&item)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := item.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (itemHandler *ItemHandler) Delete(ctx echo.Context) error {
	return ctx.NoContent(errors.ResolveErrorToCode(errors.ErrDbBadOperation))
}
