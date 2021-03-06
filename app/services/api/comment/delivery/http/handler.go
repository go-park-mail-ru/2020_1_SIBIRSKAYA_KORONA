package http

import (
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

	"github.com/labstack/echo/v4"

	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type CommentHandler struct {
	UseCase comment.UseCase
}

func CreateHandler(router *echo.Echo, useCase comment.UseCase, mw *middleware.Middleware) {
	handler := &CommentHandler{UseCase: useCase}
	router.POST("/api/boards/:bid/columns/:cid/tasks/:tid/comments", handler.Create, mw.Sanitize, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.SendSignal, mw.SendNotification, mw.CSRFmiddle)
	router.GET("/api/boards/:bid/columns/:cid/tasks/:tid/comments", handler.Get, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CSRFmiddle)
	router.DELETE("/api/boards/:bid/columns/:cid/tasks/:tid/comments/:comid", handler.Delete, mw.CheckAuth,
		mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckCommentInTask, mw.SendSignal, mw.CSRFmiddle)
}

func (commentHandler *CommentHandler) Create(ctx echo.Context) error {
	var cmt models.Comment
	body := ctx.Get("body").([]byte)
	err := cmt.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	cmt.Uid = ctx.Get("uid").(uint)
	cmt.Tid = ctx.Get("tid").(uint)
	cmt.CreatedAt = time.Now().Unix()
	err = commentHandler.UseCase.CreateComment(&cmt)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := cmt.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	// for notifications middlware
	ctx.Set("eventType", "UpdateTask")
	ctx.Set("eventType2", "AddComment")
	ctx.Set("commentText", cmt.Text)
	return ctx.String(http.StatusOK, string(resp))
}

func (commentHandler *CommentHandler) Get(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	tid := ctx.Get("tid").(uint)
	cmts, err := commentHandler.UseCase.GetComments(tid, uid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := cmts.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (commentHandler *CommentHandler) Delete(ctx echo.Context) error {
	comid := ctx.Get("comid").(uint)
	err := commentHandler.UseCase.Delete(comid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// for signal middlware
	ctx.Set("eventType", "UpdateTask")
	return ctx.NoContent(http.StatusOK)
}
