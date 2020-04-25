package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	handler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session/dilivery/grpc"
	repo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session/repository"
	useCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session/usecase"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"

	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
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
	grpcSrv := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}),
		grpc.UnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	)
	proto.RegisterSessionServer(grpcSrv, handler.CreateHandler(sesUseCase))
	log.Println("server start on address:", server.GetAddr())
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
