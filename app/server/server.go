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

	drelloMiddleware "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"
	// Нереализованные пока нами вещи заменим на стандартные из echo
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
)

type Server struct {
	IP   string
	Port uint
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.IP, server.Port)
}

func (server *Server) Run() {
	router := echo.New()

	router.Use(echoMiddleware.Logger())
	mw := drelloMiddleware.InitMiddleware()

	router.Use(mw.CORS)
	// repo
	dsn := `host=localhost user=drello_user password=drello1234 dbname=drello_db sslmode=disable`
	db, err := gorm.Open("postgres", dsn)
    if err != nil {
		log.Println("aaa")
    	log.Fatal(err)
    }
	defer db.Close()
	usrRepo := userRepo.CreateRepository(db)
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
