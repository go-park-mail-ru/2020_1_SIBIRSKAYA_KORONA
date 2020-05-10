package http

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/labstack/echo/v4"
)

type BoardHandler struct {
	useCase board.UseCase
}

func CreateHandler(router *echo.Echo, useCase board.UseCase, mw *middleware.Middleware) {
	handler := &BoardHandler{
		useCase: useCase,
	}
	// TODO: админы
	router.POST("/boards", handler.Create, mw.Sanitize, mw.CheckAuth)
	router.GET("/boards", handler.GetBoardsByUser, mw.CheckAuth)
	router.GET("/boards/:bid", handler.Get, mw.CheckAuth)
	router.GET("/boards/:bid/labels", handler.GetLabels, mw.CheckAuth, mw.CheckBoardMemberPermission)
	router.GET("/boards/:bid/columns", handler.GetColumns, mw.CheckAuth, mw.CheckBoardMemberPermission)
	router.PUT("/boards/:bid", handler.Update, mw.Sanitize, mw.CheckAuth, mw.CheckBoardAdminPermission)
	router.DELETE("/boards/:bid", handler.Delete, mw.CheckAuth, mw.CheckBoardAdminPermission)
	router.POST("/boards/:bid/members/:uid", handler.InviteMember, mw.CheckAuth, mw.CheckBoardMemberPermission, mw.SendInviteNotifications)
	router.DELETE("/boards/:bid/members/:uid", handler.DeleteMember, mw.CheckAuth, mw.CheckBoardAdminPermission)
	router.GET("/boards/:bid/search_for_invite", handler.GetUsersForInvite, mw.CheckAuth, mw.CheckBoardMemberPermission)
}

func (boardHandler *BoardHandler) Create(ctx echo.Context) error {
	var brd models.Board
	body := ctx.Get("body").([]byte)
	err := brd.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	uid := ctx.Get("uid").(uint)
	err = boardHandler.useCase.Create(uid, &brd)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := brd.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (boardHandler *BoardHandler) GetBoardsByUser(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	bAdmin, bMember, err := boardHandler.useCase.GetBoardsByUser(uid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := models.UserBoards{Admin: bAdmin, Member: bMember}.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (boardHandler *BoardHandler) Get(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	var bid uint
	_, err := fmt.Sscan(ctx.Param("bid"), &bid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	brd, err := boardHandler.useCase.Get(uid, bid, false)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := brd.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (boardHandler *BoardHandler) GetLabels(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	lbls, err := boardHandler.useCase.GetLabelsByID(bid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	body, err := lbls.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(body))
}

func (boardHandler *BoardHandler) GetColumns(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	cols, err := boardHandler.useCase.GetColumnsByID(bid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	body, err := cols.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(body))
}

func (boardHandler *BoardHandler) Update(ctx echo.Context) error {
	var brd models.Board
	body := ctx.Get("body").([]byte)
	err := brd.UnmarshalJSON(body)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	brd.ID = ctx.Get("bid").(uint)
	err = boardHandler.useCase.Update(&brd)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := brd.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (boardHandler *BoardHandler) Delete(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	err := boardHandler.useCase.Delete(bid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (boardHandler *BoardHandler) InviteMember(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	var uid uint
	_, err := fmt.Sscan(ctx.Param("uid"), &uid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err = boardHandler.useCase.InviteMember(bid, uid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	// for notifications middlware
	ctx.Set("forUid", uid)
	ctx.Set("eventType", "InviteToBoard")
	return ctx.NoContent(http.StatusOK)
}

func (boardHandler *BoardHandler) DeleteMember(ctx echo.Context) error {
	bid := ctx.Get("bid").(uint)
	var uid uint
	_, err := fmt.Sscan(ctx.Param("uid"), &uid)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	err = boardHandler.useCase.DeleteMember(bid, uid)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (boardHandler *BoardHandler) GetUsersForInvite(ctx echo.Context) error {
	nicknamePart := ctx.QueryParam("nickname")
	if nicknamePart == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}
	bid := ctx.Get("bid").(uint)
	var limit uint
	_, err := fmt.Sscan(ctx.QueryParam("limit"), &limit)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	usr, err := boardHandler.useCase.GetUsersForInvite(bid, nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	resp, err := usr.MarshalJSON()
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}
