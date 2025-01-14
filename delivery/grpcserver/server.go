package grpcserver

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/protobufmapper"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/presence"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	config Config
	srv    *presenceservice.Service
}

type Config struct {
	Port int `koanf:"port"`
}

func (s Server) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {
	res, err := s.srv.GetPresence(ctx, protobufmapper.MapFromProtobufPresenceRequestToParam(req))
	if err != nil {

		return nil, err
	}

	return protobufmapper.MapFromParamPresenceResponseToProtobuf(*res), nil
}

func New(cfg Config, srv *presenceservice.Service) Server {
	return Server{
		config:                             cfg,
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		srv:                                srv,
	}
}

func (s Server) Start() {
	address := fmt.Sprintf(":%d", s.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	presence.RegisterPresenceServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(listener); err != nil {
		logger.L().Error("Could not serve gprc server.", zap.Error(err))
	}

}
