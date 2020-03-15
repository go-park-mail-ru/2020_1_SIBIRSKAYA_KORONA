package server

import (
	"fmt"
	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/repository"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/usecase"
	"log"

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
	// user
	usrRepo := userRepo.CreateRepository()
	usrUseCase := userUseCase.CreateUseCase(usrRepo)
	userHandler.CreateHandler(router, usrUseCase)
	// session
	sesRepo := sessionRepo.CreateRepository()
	sesUseCase := sessionUseCase.CreateUseCase(sesRepo, usrRepo)
	sessionHandler.CreateHandler(router, sesUseCase)
	// start
	if err := router.Start(server.GetAddr()); err != nil {
		log.Fatal(err)
	}
}
