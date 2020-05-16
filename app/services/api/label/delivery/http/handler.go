package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/labstack/echo/v4"
)

type LabelHandler struct {
	useCase label.UseCase
}

func CreateHandler(router *echo.Echo, useCase label.UseCase, mw *middleware.Middleware) {
	handler := &LabelHandler{
		useCase: useCase,
	}
	router.POST("/api/boards/:bid/labels", handler.Create, mw.Sanitize, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.SendSignal)
	router.GET("/api/boards/:bid/labels/:lid", handler.Get, mw.CheckAuth,
		mw.CheckBoardMemberPermission)
	router.PUT("/api/boards/:bid/labels/:lid", handler.Update, mw.Sanitize, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckLabelInBoard, mw.SendSignal)
	router.DELETE("/api/boards/:bid/labels/:lid", handler.Delete, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckLabelInBoard, mw.SendSignal)
	router.POST("/api/boards/:bid/columns/:cid/tasks/:tid/labels/:lid", handler.AddLabelOnTask, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckLabelInBoard, mw.CheckColInBoard, mw.CheckTaskInCol, mw.SendSignal)
	router.DELETE("/api/boards/:bid/columns/:cid/tasks/:tid/labels/:lid", handler.RemoveLabelFromTask, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckLabelInBoard, mw.CheckColInBoard, mw.CheckTaskInCol, mw.SendSignal)
}

func (labelHandler *LabelHandler) Create(ctx echo.Context) error {
	var lbl models.Label
	body := ctx.Get("body").([]byte)
	err := lbl.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	lbl.Bid = ctx.Get("bid").(uint)
	err = labelHandler.useCase.Create(&lbl)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := lbl.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.String(http.StatusOK, string(resp))
}

func (labelHandler *LabelHandler) Get(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	var lid uint
	_, err := fmt.Sscan(ctx.Param("lid"), &lid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	lbl, err := labelHandler.useCase.Get(bid, lid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := lbl.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (labelHandler *LabelHandler) Update(ctx echo.Context) error {
	var lbl models.Label
	body := ctx.Get("body").([]byte)
	err := lbl.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	lbl.ID = ctx.Get("lid").(uint)
	lbl.Bid = ctx.Get("bid").(uint)
	err = labelHandler.useCase.Update(lbl)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.NoContent(http.StatusOK)
}

func (labelHandler *LabelHandler) Delete(ctx echo.Context) error {
	lid := ctx.Get("lid").(uint)
	err := labelHandler.useCase.Delete(lid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.NoContent(http.StatusOK)
}

func (labelHandler *LabelHandler) AddLabelOnTask(ctx echo.Context) error {
	lid := ctx.Get("lid").(uint)
	tid := ctx.Get("tid").(uint)
	err := labelHandler.useCase.AddLabelOnTask(lid, tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.NoContent(http.StatusOK)
}

func (labelHandler *LabelHandler) RemoveLabelFromTask(ctx echo.Context) error {
	lid := ctx.Get("lid").(uint)
	tid := ctx.Get("tid").(uint)
	err := labelHandler.useCase.RemoveLabelFromTask(lid, tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.NoContent(http.StatusOK)
}
