package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/interceptor"
	handler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/delivery/grpc"
	repo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/repository"
	useCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/user/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/metric"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	IP     string
	Port   uint
	Config *config.UserConfigController
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.IP, server.Port)
}

// TODO: конфиг
func (server *Server) Run() {
	logger.InitLogger()
	// repo
	postgresClient, err := gorm.Open(server.Config.GetDB(), server.Config.GetDBConnection())
	if err != nil {
		log.Fatal(err)
	}
	defer postgresClient.Close()
	usrRepo := repo.CreateRepository(postgresClient)
	// usecase
	usrUseCase := useCase.CreateUseCase(usrRepo)
	// handler
	listener, err := net.Listen("tcp", server.GetAddr())
	if err != nil {
		log.Fatal(err)
	}
	metr, err := metric.CreateMetrics(server.Config.GetMetricsURL(), server.Config.GetServiceName())
	if err != nil {
		log.Fatal(err)
	}
	mw := interceptor.CreateInterceptor(metr)
	grpcSrv := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}),
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor(), mw.Metrics),
	)
	proto.RegisterUserServer(grpcSrv, handler.CreateHandler(usrUseCase))
	log.Println("server start on address:", server.GetAddr())
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
