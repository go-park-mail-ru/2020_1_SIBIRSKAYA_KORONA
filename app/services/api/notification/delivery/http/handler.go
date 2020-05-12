package http

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	UseCase notification.UseCase
}

func CreateHandler(router *echo.Echo, useCase notification.UseCase, mw *middleware.Middleware) {
	handler := &NotificationHandler{UseCase: useCase}

	router.GET("/notifications", handler.GetAll, mw.CheckAuth)
	router.PUT("/notifications", handler.UpdateAll, mw.CheckAuth)
	router.DELETE("/notifications", handler.DeleteAll, mw.CheckAuth)
}

func (notificationHandler *NotificationHandler) GetAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	events, has := notificationHandler.UseCase.GetAll(uid)
	if !has {
		logger.Error("no notifications for the user", uid)
	}
	resp, err := events.MarshalJSON()
	if err != nil {
		logger.Error(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, string(resp))
}

func (notificationHandler *NotificationHandler) UpdateAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	if err := notificationHandler.UseCase.UpdateAll(uid); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

func (notificationHandler *NotificationHandler) DeleteAll(ctx echo.Context) error {
	uid := ctx.Get("uid").(uint)
	if err := notificationHandler.UseCase.DeleteAll(uid); err != nil {
		logger.Error(err)
		return ctx.String(errors.ResolveErrorToCode(err), err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}
