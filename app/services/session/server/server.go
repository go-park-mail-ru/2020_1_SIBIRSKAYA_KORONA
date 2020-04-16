package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	handler "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session/dilivery/grpc"
	sessionRepo "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session/repository"
	sessionUseCase "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	IP   string
	Port uint
}

func (server *Server) GetAddr() string {
	return fmt.Sprintf("%s:%d", server.IP, server.Port)
}

// TODO: логгер, конфиг
func (server *Server) Run() {
	// repo
	memCacheClient := memcache.New("127.0.0.1:11211")
	err := memCacheClient.Ping()
	defer memCacheClient.DeleteAll()
	sesRepo := sessionRepo.CreateRepository(memCacheClient)
	// usecase
	sesUseCase := sessionUseCase.CreateUseCase(sesRepo)
	// handler
	listener, err := net.Listen("tcp", server.GetAddr())
	if err != nil {
		log.Fatal(err)
	}
	grpcSrv := grpc.NewServer(grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		},
	))
	proto.RegisterSessionServer(grpcSrv, handler.CreateHandler(sesUseCase))
	log.Println("server start on address:", server.GetAddr())
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
