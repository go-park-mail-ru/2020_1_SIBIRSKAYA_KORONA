package webSocketPool

type WebSocketPool interface {
	Send(uid uint, mess []byte)
}
