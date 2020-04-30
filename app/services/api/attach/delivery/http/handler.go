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

	router.GET("/boards/:bid/columns/:cid/tasks/:tid/files", handler.GetFiles,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	router.POST("/boards/:bid/columns/:cid/tasks/:tid/files", handler.Create,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	// mw.Sanitaze == death ?
	router.DELETE("/boards/:bid/columns/:cid/tasks/:tid/files/:fid", handler.Delete,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckAttachInTask)
}

func (attachHandler *AttachHandler) Create(ctx echo.Context) error {
	tid := ctx.Get("tid").(uint)
	attach := &models.AttachedFile{Tid: tid}

	AttachDescriptor, err := ctx.FormFile("file")
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	logger.Info("File: ", AttachDescriptor.Filename)

	err = attachHandler.UseCase.Create(attach, AttachDescriptor)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := attach.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
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

	return ctx.NoContent(http.StatusOK)
}
