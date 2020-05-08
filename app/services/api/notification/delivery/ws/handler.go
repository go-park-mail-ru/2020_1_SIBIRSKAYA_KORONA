package ws

import (
	"fmt"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	UseCase notification.UseCase
	upgrader websocket.Upgrader
}

func CreateHandler(router *echo.Echo, useCase notification.UseCase, mw *middleware.Middleware) {
	handler := &NotificationHandler{
		UseCase: useCase,
		//upgrader: websocket.Upgrader{},
	}
	router.GET("/ws", handler.Run, mw.CheckAuth)
}

func (notificationHandler *NotificationHandler) Run(ctx echo.Context) error {
	ws, err := notificationHandler.upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	uid := ctx.Get("id")
	for {
		// Write
		events, has := notificationHandler.UseCase.GetEvents(uid)
		if has {
			events.
			err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
			if err != nil {
				logger.Error(err)
			}
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			logger.Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}