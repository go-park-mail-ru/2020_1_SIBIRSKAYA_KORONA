package server

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/metric"

	labelHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/delivery/http"

	notificationHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/delivery/ws"

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

	commentHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/delivery/http"
	commentRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/repository"
	commentUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/comment/usecase"

	checklistHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist/delivery/http"
	checklistRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist/repository"
	checklistUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/checklist/usecase"

	itemHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/delivery/http"
	itemRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/repository"
	itemUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/item/usecase"

	attachHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/delivery/http"
	attachRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/repository"
	attachUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/attach/usecase"

	labelRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/repository"
	labelUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/usecase"

	drelloMiddleware "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
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
	// session
	grpcSessionConn, err := grpc.Dial(
		server.ApiConfig.GetSessionClient(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcSessionConn.Close()
	sessionGrpcClient := proto.NewSessionClient(grpcSessionConn)
	// user
	grpcUserConn, err := grpc.Dial(
		server.ApiConfig.GetUserClient(),
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
	}
	logger.Info("Postgresql succesfull start")
	defer postgresClient.Close()
	postgresClient.AutoMigrate(&models.User{}, &models.Board{}, &models.Column{}, &models.Task{}, &models.Comment{},
		&models.Checklist{}, &models.Item{}, &models.Label{}, &models.AttachedFile{}, &models.Event{})
	sesRepo := sessionRepo.CreateRepository(sessionGrpcClient)
	usrRepo := userRepo.CreateRepository(userGrpcClient, server.UserConfig)
	brdRepo := boardRepo.CreateRepository(postgresClient)
	colRepo := colsRepo.CreateRepository(postgresClient)
	lblRepo := labelRepo.CreateRepository(postgresClient)
	tskRepo := taskRepo.CreateRepository(postgresClient)
	comRepo := commentRepo.CreateRepository(postgresClient)
	chlistRepo := checklistRepo.CreateRepository(postgresClient)
	itmRepo := itemRepo.CreateRepository(postgresClient)

	S3session, err := session.NewSession(
		&aws.Config{Region: aws.String(server.ApiConfig.GetS3BucketRegion())},
	)
	if err != nil {
		logger.Fatal(err)
	}

	attachModelRepo := attachRepo.CreateRepository(postgresClient)
	attachFileRepo := attachRepo.CreateS3Repository(S3session, server.ApiConfig.GetS3Bucket())

	// use case
	sUseCase := sessionUseCase.CreateUseCase(sesRepo, usrRepo)
	uUseCase := userUseCase.CreateUseCase(sesRepo, usrRepo)
	bUseCase := boardUseCase.CreateUseCase(usrRepo, brdRepo)
	cUseCase := colsUseCase.CreateUseCase(colRepo)
	lUseCase := labelUseCase.CreateUseCase(lblRepo)
	tUseCase := taskUseCase.CreateUseCase(tskRepo, usrRepo)
	comUseCase := commentUseCase.CreateUseCase(comRepo)
	chUseCase := checklistUseCase.CreateUseCase(chlistRepo, itmRepo)
	itmUseCase := itemUseCase.CreateUseCase(itmRepo)
	atchUseCase := attachUseCase.CreateUseCase(attachModelRepo, attachFileRepo)

	// middlware
	router := echo.New()
	metr, err := metric.CreateMetrics(server.ApiConfig.GetMetricsURL(), server.ApiConfig.GetServiceName())
	if err != nil {
		log.Fatal(err)
	}
	mw := drelloMiddleware.CreateMiddleware(metr, sUseCase, bUseCase, cUseCase, tUseCase, comUseCase, chUseCase, itmUseCase, lUseCase, atchUseCase)
	router.Use(mw.RequestLogger)
	router.Use(mw.CORS)
	router.Use(mw.ProcessPanic)
	router.Use(mw.Metrics)
	// delivery
	sessionHandler.CreateHandler(router, sUseCase, mw)
	userHandler.CreateHandler(router, uUseCase, mw)
	boardHandler.CreateHandler(router, bUseCase, mw)
	colsHandler.CreateHandler(router, cUseCase, mw)
	labelHandler.CreateHandler(router, lUseCase, mw)
	taskHandler.CreateHandler(router, tUseCase, mw)
	commentHandler.CreateHandler(router, comUseCase, mw)
	checklistHandler.CreateHandler(router, chUseCase, mw)
	itemHandler.CreateHandler(router, itmUseCase, mw)
	attachHandler.CreateHandler(router, atchUseCase, mw)

	notificationHandler.CreateHandler(router, nil, mw)

	// start
	if err := router.StartTLS(server.GetAddr(), server.ApiConfig.GetTLSCrtPath(), server.ApiConfig.GetTLSKeyPath()); err != nil {
		log.Fatal(err)
	}
}
