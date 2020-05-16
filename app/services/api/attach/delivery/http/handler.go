package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/labstack/echo/v4"
)

type AttachHandler struct {
	UseCase attach.UseCase
}

func CreateHandler(router *echo.Echo, useCase attach.UseCase, mw *middleware.Middleware) {
	handler := &AttachHandler{
		UseCase: useCase,
	}

	router.POST("/api/boards/:bid/columns/:cid/tasks/:tid/files", handler.Create, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.SendSignal)
	router.GET("/api/boards/:bid/columns/:cid/tasks/:tid/files", handler.GetFiles, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	// mw.Sanitaze == death ?
	router.DELETE("/api/boards/:bid/columns/:cid/tasks/:tid/files/:fid", handler.Delete, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckAttachInTask, mw.SendSignal)
}

func (attachHandler *AttachHandler) Create(ctx echo.Context) error {
	tid := ctx.Get("tid").(uint)
	atch := &models.AttachedFile{Tid: tid}
	AttachDescriptor, err := ctx.FormFile("file")
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}
	logger.Info("File: ", AttachDescriptor.Filename)
	err = attachHandler.UseCase.Create(atch, AttachDescriptor)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := atch.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.String(http.StatusOK, string(resp))

}

func (attachHandler *AttachHandler) GetFiles(ctx echo.Context) error {
	tid := ctx.Get("tid").(uint)
	attachedFiles, err := attachHandler.UseCase.Get(tid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := attachedFiles.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))

}

func (attachHandler *AttachHandler) Delete(ctx echo.Context) error {
	fid := ctx.Get("fid").(uint)
	err := attachHandler.UseCase.Delete(fid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.NoContent(http.StatusOK)
}
