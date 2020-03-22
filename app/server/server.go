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
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/bradfitz/gomemcache/memcache"
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

	router.Use(echoMiddleware.Logger())
	mw := drelloMiddleware.InitMiddleware()

	router.Use(mw.CORS)
	// router.Use(mw.ProcessPanic)
	// repo
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
		log.Fatal(err)
	}
	defer postgresClient.Close()
	usrRepo := userRepo.CreateRepository(postgresClient)
	// memCache
	memcacheHost := viper.GetString("memcached.host")
	memcachePort := viper.GetString("memcached.port")
	memcacheConnection := fmt.Sprintf("%s:%s", memcacheHost, memcachePort)
	memCacheClient := memcache.New(memcacheConnection)
	err = memCacheClient.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer memCacheClient.DeleteAll()
	sesRepo := sessionRepo.CreateRepository(memCacheClient)
	// use case
	sesUseCase := sessionUseCase.CreateUseCase(sesRepo, usrRepo)
	usrUseCase := userUseCase.CreateUseCase(sesRepo, usrRepo)
	// delivery
	userHandler.CreateHandler(router, usrUseCase, mw)
	sessionHandler.CreateHandler(router, sesUseCase, mw)
	// start
	if err := router.Start(server.GetAddr()); err != nil {
		log.Fatal(err)
	}
}
