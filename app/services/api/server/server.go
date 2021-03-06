package server

import (
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/spf13/viper"

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

	labelHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/delivery/http"
	labelRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/repository"
	labelUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/label/usecase"

	templateHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/template/delivery/http"
	templateUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/template/usecase"

	notificationHandler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/delivery/http"
	notificationRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/repository"
	notificationUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/notification/usecase"

	drelloMiddleware "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/middleware"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/metric"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/webSocketPool/gorillaWs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

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
		&models.Checklist{}, &models.Item{}, &models.Label{}, &models.AttachedFile{}, &models.Event{}, &models.EventMetaData{})
	sesRepo := sessionRepo.CreateRepository(sessionGrpcClient)
	usrRepo := userRepo.CreateRepository(userGrpcClient, server.UserConfig)
	brdRepo := boardRepo.CreateRepository(postgresClient)
	colRepo := colsRepo.CreateRepository(postgresClient)
	lblRepo := labelRepo.CreateRepository(postgresClient)
	tskRepo := taskRepo.CreateRepository(postgresClient)
	comRepo := commentRepo.CreateRepository(postgresClient)
	chlistRepo := checklistRepo.CreateRepository(postgresClient)
	itmRepo := itemRepo.CreateRepository(postgresClient)
	ntftRepo := notificationRepo.CreateRepository(postgresClient)

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
	ntftUseCase := notificationUseCase.CreateUseCase(usrRepo, ntftRepo)

	templateReadersMap := ReadTemplates(server.ApiConfig.GetTemplatesPath())
	tmplUseCase := templateUseCase.CreateUseCase(lblRepo, colRepo, tskRepo, comRepo, chlistRepo, brdRepo, usrRepo, templateReadersMap)

	// middlware
	router := echo.New()

	mw := drelloMiddleware.CreateMiddleware(sUseCase, uUseCase, bUseCase, cUseCase, tUseCase,
		comUseCase, chUseCase, itmUseCase, lUseCase, atchUseCase, ntftUseCase)

	router.Use(mw.RequestLogger)
	router.Use(mw.CORS)
	router.Use(mw.ProcessPanic)
	router.Use(mw.Metrics)

	// metrics
	metr, err := metric.CreateMetrics(server.ApiConfig.GetMetricsURL(), server.ApiConfig.GetServiceName())
	if err != nil {
		log.Fatal(err)
	}
	mw.SetMetrics(metr)

	// wsPool
	wsPool := gorillaWs.CreateWebSocketPool(router, mw)
	mw.SetWebSocketPool(wsPool)

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
	notificationHandler.CreateHandler(router, ntftUseCase, mw)
	templateHandler.CreateHandler(router, tmplUseCase, mw)

	// start
	err = router.StartTLS(server.GetAddr(), server.ApiConfig.GetTLSCrtPath(), server.ApiConfig.GetTLSKeyPath())
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: Вынести в более подходящее место
func ReadTemplates(tmplPath string) map[string]*viper.Viper {
	// Временный костыль - один независимый экземляр viper на каждый шаблон
	templateReadersMap := make(map[string]*viper.Viper)

	weekTemplateReader := viper.New()
	weekTemplateReader.AddConfigPath(tmplPath)
	weekTemplateReader.SetConfigName("week_plan")
	err := weekTemplateReader.MergeInConfig()
	if err != nil {
		logger.Fatal(err)
	}
	templateReadersMap["week_plan"] = weekTemplateReader

	projectTemplateReader := viper.New()
	projectTemplateReader.AddConfigPath(tmplPath)
	projectTemplateReader.SetConfigName("product_management")
	err = projectTemplateReader.MergeInConfig()
	if err != nil {
		logger.Fatal(err)
	}
	templateReadersMap["product_management"] = projectTemplateReader

	return templateReadersMap
}
