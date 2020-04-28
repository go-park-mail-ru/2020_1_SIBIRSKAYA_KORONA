package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/interceptor"
	handler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session/delivery/grpc"
	repo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session/repository"
	useCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/metric"

	"github.com/bradfitz/gomemcache/memcache"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	IP     string
	Port   uint
	Config *config.SessionConfigController
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.IP, server.Port)
}

func (server *Server) Run() {
	// repo
	memCacheClient := memcache.New(server.Config.GetMemcachedConnect())
	err := memCacheClient.Ping()
	defer memCacheClient.DeleteAll()
	sesRepo := repo.CreateRepository(memCacheClient)
	// usecase
	sesUseCase := useCase.CreateUseCase(sesRepo)
	// handler
	listener, err := net.Listen("tcp", server.GetAddr())
	if err != nil {
		log.Fatal(err)
	}
	metr, err := metric.CreateMetrics("0.0.0.0:7071", "session")
	if err != nil {
		log.Fatal(err)
	}
	mw := interceptor.CreateInterceptor(metr)
	grpcSrv := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}),
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor(), mw.Metrics),
	)
	proto.RegisterSessionServer(grpcSrv, handler.CreateHandler(sesUseCase))
	log.Println("server start on address:", server.GetAddr())
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
