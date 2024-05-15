package service

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"

	"github.com/lakhansamani/container-orchestrator-apis/container"
	mp "github.com/lakhansamani/container-orchestrator/internal/memorystore/providers"
)

type Dependencies struct {
	Logger       zerolog.Logger
	DockerClient *client.Client
	MemoryStore  mp.MemoryStoreProvider
}

// Service implements the Data service.
type Service interface {
	container.ContainerServiceServer

	// Run the service until the given context is canceled
	Run(context.Context) error
}

type service struct {
	Dependencies
}

// New creates a new Service.
func New(ctx context.Context, deps Dependencies) (Service, error) {
	return &service{
		Dependencies: deps,
	}, nil
}

// Run this service
func (s *service) Run(ctx context.Context) error {
	return nil
}
