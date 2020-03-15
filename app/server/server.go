package server

import (
	"fmt"
	"log"

	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/repository"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/usecase"

	sessionHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/delivery/http"
	sessionRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/repository"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/usecase"

	"github.com/labstack/echo/v4"
)

type Server struct {
	IP string
	Port uint
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.IP, server.Port)
}

func (server *Server) Run() {
	log.Println("init")
	router := echo.New()
	// repo
	usrRepo := userRepo.CreateRepository()
	sesRepo := sessionRepo.CreateRepository()
	// use case
	sesUseCase := sessionUseCase.CreateUseCase(sesRepo, usrRepo)
	usrUseCase := userUseCase.CreateUseCase(sesRepo, usrRepo)
	// delivery
	userHandler.CreateHandler(router, usrUseCase)
	sessionHandler.CreateHandler(router, sesUseCase)
	// start
	if err := router.Start(server.GetAddr()); err != nil {
		log.Fatal(err)
	}
}
