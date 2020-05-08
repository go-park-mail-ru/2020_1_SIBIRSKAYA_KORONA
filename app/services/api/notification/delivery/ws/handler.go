package ws

import (
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	UseCase  notification.UseCase
	upgrader websocket.Upgrader
}

func CreateHandler(router *echo.Echo, useCase notification.UseCase, mw *middleware.Middleware) {
	handler := &NotificationHandler{
		UseCase:  useCase,
		upgrader: websocket.Upgrader{},
	}
	router.GET("ws", handler.Run, mw.CheckAuth)
}

func (notificationHandler *NotificationHandler) Run(ctx echo.Context) error {
	ws, err := notificationHandler.upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	uid := ctx.Get("uid").(uint)
	for {
		time.Sleep(10 * time.Second)
		events, has := notificationHandler.UseCase.Pop(uid)
		if !has {
			continue
		}
		res, err := events.MarshalJSON()
		if err != nil {
			logger.Error(err)
			continue
		}
		err = ws.WriteMessage(websocket.TextMessage, res)
		if err != nil {
			logger.Error(err)
		}
	}
}
