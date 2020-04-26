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
	useCase comment.UseCase
}

func CreateHandler(router *echo.Echo, useCase comment.UseCase, mw *middleware.Middleware) {
	handler := &CommentHandler{useCase: useCase}

	router.POST("boards/:bid/columns/:cid/tasks/:tid/comments", handler.CreateComment,
		mw.Sanitize, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	router.GET("boards/:bid/columns/:cid/tasks/:tid/comments", handler.GetComments,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol)
	router.DELETE("boards/:bid/columns/:cid/tasks/:tid/comments/:comid", handler.DeleteComment,
		mw.CheckAuth, mw.CheckBoardMemberPermission, mw.CheckColInBoard, mw.CheckTaskInCol, mw.CheckCommentInTask)
	//TODO: mw.CheckCommentInTask ?

}

func (commentHandler *CommentHandler) CreateComment(ctx echo.Context) error {
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

	err = commentHandler.useCase.CreateComment(&cmt)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := cmt.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (commentHandler *CommentHandler) GetComments(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	tid := ctx.Get("tid").(uint)

	cmts, err := commentHandler.useCase.GetComments(tid, uid)
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

func (commentHandler *CommentHandler) DeleteComment(ctx echo.Context) error {
	fid := ctx.Get("comid").(uint)

	err := commentHandler.useCase.Delete(fid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
