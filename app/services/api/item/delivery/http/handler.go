package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/labstack/echo/v4"
)

type ItemHandler struct {
	UseCase item.UseCase
}

func CreateHandler(router *echo.Echo, useCase item.UseCase, mw *middleware.Middleware) {
	handler := &ItemHandler{
		UseCase: useCase,
	}
	router.POST("/api/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid/items", handler.Create,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard,
		mw.CheckTaskInCol, mw.CheckChecklistInTask, mw.SendSignal)
	router.PUT("/api/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid/items/:itid", handler.Update,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol,
		mw.CheckChecklistInTask, mw.CheckItemInChecklist, mw.SendSignal)
	router.DELETE("/api/boards/:bid/columns/:cid/tasks/:tid/checklists/:clid/items/:itid", handler.Delete,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol,
		mw.CheckChecklistInTask, mw.CheckItemInChecklist)
}

func (itemHandler *ItemHandler) Create(ctx echo.Context) error {
	clid := ctx.Get("clid").(uint)
	var itm models.Item
	body := ctx.Get("body").([]byte)
	err := itm.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	itm.Clid = clid
	err = itemHandler.UseCase.Create(&itm)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := itm.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.String(http.StatusOK, string(resp))
}

func (itemHandler *ItemHandler) Update(ctx echo.Context) error {
	clid := ctx.Get("clid").(uint)
	itid := ctx.Get("itid").(uint)
	var itm models.Item
	body := ctx.Get("body").([]byte)
	err := itm.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	itm.Clid = clid
	itm.ID = itid
	err = itemHandler.UseCase.Update(&itm)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := itm.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.String(http.StatusOK, string(resp))
}

func (itemHandler *ItemHandler) Delete(ctx echo.Context) error {
	return ctx.NoContent(errors.ResolveErrorToCode(errors.ErrDbBadOperation))
}
