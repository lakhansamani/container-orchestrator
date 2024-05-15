package cmd

import (
	"context"
	"os"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/lakhansamani/container-orchestrator/internal/memorystore"
	"github.com/lakhansamani/container-orchestrator/internal/server"
	"github.com/lakhansamani/container-orchestrator/internal/service"
)

var (
	// RootCmd is the root (and only) command of this service
	RootCmd = &cobra.Command{
		Use:   "container-orchestrator",
		Short: "Container Orchestrator",
		Run:   runRootCmd,
	}
	rootArgs struct {
		version  string
		logLevel string
		redisURL string
		server   struct {
			Host     string
			GRPCPort int
		}
	}
)

// SetVersion stores the given version
func SetVersion(version string) {
	rootArgs.version = version
}

func init() {
	f := RootCmd.Flags()
	// Logging flags
	f.StringVar(&rootArgs.logLevel, "log-level", "", "Minimum log level")
	// Server flags
	f.StringVar(&rootArgs.server.Host, "host", "0.0.0.0", "Host interface to listen on")
	f.IntVar(&rootArgs.server.GRPCPort, "grpc-port", 5600, "Port to listen on for GRPC requests")
	f.StringVar(&rootArgs.redisURL, "redis-url", "redis://localhost:6379", "URL of the Redis server")
}

// Run the service
func runRootCmd(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	dockerCLI, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create docker client")
	}
	memoryStoreProvider, err := memorystore.NewMemoryStore(rootArgs.redisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create memory store")
	}
	svc, err := service.New(ctx, service.Dependencies{
		Logger:       log,
		DockerClient: dockerCLI,
		MemoryStore:  memoryStoreProvider,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create service")
	}

	srv, err := server.New(server.Config{
		Host:     rootArgs.server.Host,
		GRPCPort: rootArgs.server.GRPCPort,
	}, log, svc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create server")
	}
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return srv.Run(ctx) })
	g.Go(func() error { return svc.Run(ctx) })
	if err := g.Wait(); err != nil && err != context.Canceled {
		log.Fatal().Err(err).Msg("Failed to run")
	}
}
