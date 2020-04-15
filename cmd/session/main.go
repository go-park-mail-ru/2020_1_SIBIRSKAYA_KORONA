package main

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session/server"
)

func main() {
	srv := &server.Server{
		IP:   "127.0.0.1",
		Port: 8081,
	}
	srv.Run()
}
