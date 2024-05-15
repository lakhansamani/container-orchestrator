package server

import (
	"context"
	"net"
	"strconv"

	"github.com/lakhansamani/container-orchestrator-apis/container"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Config is the configuration for the server
type Config struct {
	Host     string
	GRPCPort int
}

// Server is the server
type Server struct {
	Config
	service Service
	log     zerolog.Logger
}

// Service is the service
type Service interface {
	container.ContainerServiceServer
}

// New configures a new Server.
func New(cfg Config, log zerolog.Logger, service Service) (*Server, error) {
	return &Server{
		Config:  cfg,
		service: service,
		log:     log,
	}, nil
}

// Run the server until the given context is canceled.
func (s *Server) Run(ctx context.Context) error {
	// Prepare GRPC listener
	log := s.log
	grpcAddr := net.JoinHostPort(s.Host, strconv.Itoa(s.GRPCPort))
	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listen on address %s", grpcAddr)
	}
	log.Debug().Msgf("listening on %s", grpcAddr)
	grpcServer := grpc.NewServer()
	container.RegisterContainerServiceServer(grpcServer, s.service)
	if err := grpcServer.Serve(grpcLis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve grpc")
	}
	return nil
}
