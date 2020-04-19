package server

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"

	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/repository"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user/usecase"

	sessionHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/delivery/http"
	sessionRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/repository"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/session/usecase"

	boardHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board/delivery/http"
	boardRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board/repository"
	boardUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/board/usecase"

	colsHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/delivery/http"
	colsRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/repository"
	colsUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/column/usecase"

	taskHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/delivery/http"
	taskRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/repository"
	taskUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/task/usecase"

	drelloMiddleware "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type Server struct {
	IP         string
	Port       uint
	ApiConfig  *config.ApiConfigController
	UserConfig *config.UserConfigController
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.IP, server.Port)
}

func (server *Server) Run() {
	// repo
	// micro-serv
	// TODO: конфиг
	// session
	grpcSessionConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcSessionConn.Close()
	sessionGrpcClient := proto.NewSessionClient(grpcSessionConn)
	// user
	grpcUserConn, err := grpc.Dial(
		"127.0.0.1:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcUserConn.Close()
	userGrpcClient := proto.NewUserClient(grpcUserConn)
	// postgres
	postgresClient, err := gorm.Open(server.ApiConfig.GetDB(), server.ApiConfig.GetDBConnection())
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("Postgresql succesfull start")
	}
	defer postgresClient.Close()
	postgresClient.AutoMigrate(&models.User{}, &models.Board{}, &models.Column{}, &models.Task{})
	sesRepo := sessionRepo.CreateRepository(sessionGrpcClient)
	usrRepo := userRepo.CreateRepository(userGrpcClient, server.UserConfig)
	brdRepo := boardRepo.CreateRepository(postgresClient)
	colRepo := colsRepo.CreateRepository(postgresClient)
	tskRepo := taskRepo.CreateRepository(postgresClient)

	// use case
	sUseCase := sessionUseCase.CreateUseCase(sesRepo, usrRepo)
	uUseCase := userUseCase.CreateUseCase(sesRepo, usrRepo)
	bUseCase := boardUseCase.CreateUseCase(usrRepo, brdRepo)
	cUseCase := colsUseCase.CreateUseCase(colRepo)
	tUseCase := taskUseCase.CreateUseCase(tskRepo, usrRepo)

	// delivery
	mw := drelloMiddleware.CreateMiddleware(sUseCase, bUseCase, cUseCase, tUseCase)
	router := echo.New()
	router.Use(mw.RequestLogger)
	router.Use(mw.CORS)
	router.Use(mw.ProcessPanic)
	sessionHandler.CreateHandler(router, sUseCase, mw)
	userHandler.CreateHandler(router, uUseCase, mw)
	boardHandler.CreateHandler(router, bUseCase, mw)
	colsHandler.CreateHandler(router, cUseCase, mw)
	taskHandler.CreateHandler(router, tUseCase, mw)

	// start
	if err := router.Start(server.GetAddr()); err != nil {
		log.Fatal(err)
	}
}
