package gorillaWs

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/webSocketPool"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type GorillaWebSocketPool struct {
	socketPool map[uint][]*websocket.Conn
	mux        *sync.Mutex
	upgrader   websocket.Upgrader
}

func CreateWebSocketPool(router *echo.Echo, mw *middleware.Middleware) webSocketPool.WebSocketPool {
	wsPool := &GorillaWebSocketPool{
		socketPool: make(map[uint][]*websocket.Conn),
		mux:        &sync.Mutex{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	router.GET("/api/ws", wsPool.run, mw.CheckAuth)
	return wsPool
}

func (wsPool *GorillaWebSocketPool) run(ctx echo.Context) error {
	ws, err := wsPool.upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer ws.Close()
	uid := ctx.Get("uid").(uint)
	err = wsPool.Add(uid, ws)
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Info("add web socket to pool for user:", uid)

	defer func() {
		errDelete := wsPool.Delete(uid, ws)
		if errDelete != nil {
			log.Fatal(errDelete)
		}
	}()

	var readErr error = nil
	for readErr == nil {
		_, _, readErr = ws.ReadMessage()
		logger.Error(readErr)
	}
	logger.Info("close web socket with error:", readErr)
	return ctx.NoContent(http.StatusOK)
}

func (wsPool *GorillaWebSocketPool) Add(uid uint, ws *websocket.Conn) error {
	wsPool.mux.Lock()
	defer wsPool.mux.Unlock()
	for _, elem := range wsPool.socketPool[uid] {
		if elem == ws {
			return errors.New("conflict")
		}
	}
	wsPool.socketPool[uid] = append(wsPool.socketPool[uid], ws)
	return nil
}

func (wsPool *GorillaWebSocketPool) Delete(uid uint, ws *websocket.Conn) error {
	wsPool.mux.Lock()
	defer wsPool.mux.Unlock()
	for idx, elem := range wsPool.socketPool[uid] {
		if elem == ws {
			wsPool.socketPool[uid][idx] = wsPool.socketPool[uid][len(wsPool.socketPool[uid])-1]
			wsPool.socketPool[uid] = wsPool.socketPool[uid][:len(wsPool.socketPool[uid])-1]
			logger.Info("close web socket connection for user:", uid)
			return nil
		}
	}
	return errors.New("not found ws")
}

func (wsPool *GorillaWebSocketPool) Send(uid uint, mess []byte) {
	wsPool.mux.Lock()
	defer wsPool.mux.Unlock()
	for _, elem := range wsPool.socketPool[uid] {
		err := elem.WriteMessage(websocket.TextMessage, mess)
		if err != nil {
			logger.Error("error while sending by ws to user:", uid, " error:", err)
		}
	}
}
