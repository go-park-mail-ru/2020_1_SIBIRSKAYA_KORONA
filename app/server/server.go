package server

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

	userHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/delivery/http"
	userRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/repository"
	userUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user/usecase"

	sessionHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/delivery/http"
	sessionRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/repository"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/usecase"

	boardHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board/delivery/http"
	boardRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board/repository"
	boardUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/board/usecase"

	drelloMiddleware "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/middleware"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
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

	mw := drelloMiddleware.InitMiddleware()

	router.Use(mw.CORS)
	router.Use(mw.ProcessPanic)

	// repo

	// memCache
	memCacheHost := viper.GetString("memcached.host")
	memCachePort := viper.GetString("memcached.port")
	memCacheConnection := fmt.Sprintf("%s:%s", memCacheHost, memCachePort)
	memCacheClient := memcache.New(memCacheConnection)
	err := memCacheClient.Ping()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("Memcached succesfull start")
	}
	defer memCacheClient.DeleteAll()
	sesRepo := sessionRepo.CreateRepository(memCacheClient)

	// postgres
	dbms := viper.GetString("database.dbms")
	dbHost := viper.GetString("database.host")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	dbMode := viper.GetString("database.sslmode")
	dbConnection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbMode)
	postgresClient, err := gorm.Open(dbms, dbConnection)
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Info("Postgresql succesfull start")
	}
	postgresClient.AutoMigrate(&models.User{}, &models.Board{})
	defer postgresClient.Close()
	usrRepo := userRepo.CreateRepository(postgresClient)
	bRepo := boardRepo.CreateRepository(postgresClient)

	// use case
	sUseCase := sessionUseCase.CreateUseCase(sesRepo, usrRepo)
	uUseCase := userUseCase.CreateUseCase(sesRepo, usrRepo)
	bUseCase := boardUseCase.CreateUseCase(sesRepo, usrRepo, bRepo)
	// delivery
	sessionHandler.CreateHandler(router, sUseCase, mw)
	userHandler.CreateHandler(router, uUseCase, mw)
	boardHandler.CreateHandler(router, bUseCase, mw)

	// start
	if err := router.Start(server.GetAddr()); err != nil {
		log.Fatal(err)
	}
}
