package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	dc "github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lakhansamani/container-orchestrator-apis/container"
)

// CreateContainer creates a new container
func (s *service) CreateContainer(ctx context.Context, req *container.CreateContainerRequest) (*container.Container, error) {
	envs := []string{}
	for _, env := range req.GetEnvVars() {
		envs = append(envs, env.GetKey()+"="+env.GetValue())
	}
	cfg := &dc.Config{
		Image: req.GetImage(),
		Env:   envs,
	}
	reader, err := s.DockerClient.ImagePull(ctx, req.GetImage(), types.ImagePullOptions{})
	if err != nil {
		log.Debug().Err(err).Msg("Error pulling image")
		return nil, status.Error(codes.Internal, err.Error())
	}
	io.Copy(os.Stdout, reader)

	res, err := s.DockerClient.ContainerCreate(ctx, cfg, nil, nil, nil, req.GetName())
	if err != nil {
		log.Debug().Err(err).Msg("Error creating container")
		return nil, status.Error(codes.Internal, err.Error())
	}
	memoryStoreKey := fmt.Sprintf("%s:%s", res.ID, req.GetName())
	s.MemoryStore.SetData(memoryStoreKey, "created")
	go func() {
		ctx := context.Background()
		if err := s.DockerClient.ContainerStart(ctx, res.ID, dc.StartOptions{}); err != nil {
			s.MemoryStore.SetData(memoryStoreKey, fmt.Sprintf("failed: %s", err.Error()))
		} else {
			log.Debug().Str("container_id", res.ID).Msgf("Container started")
			s.MemoryStore.SetData(memoryStoreKey, "started")
			// Wait for the container to running or exited with no error (0)
			// If the container exits with an error, update the memory store with the error message
			for {
				containerJSON, err := s.DockerClient.ContainerInspect(ctx, res.ID)
				log.Debug().Interface("container", containerJSON.State).Msg("Container status")
				if err != nil {
					s.MemoryStore.SetData(memoryStoreKey, fmt.Sprintf("failed: %s", err.Error()))
					break
				}
				if containerJSON.State.Status == "running" || containerJSON.State.Status == "exited" && containerJSON.State.ExitCode == 0 {
					s.MemoryStore.SetData(memoryStoreKey, containerJSON.State.Status)
					break
				} else if containerJSON.State.Status == "exited" {
					s.MemoryStore.SetData(memoryStoreKey, fmt.Sprintf("failed: %s, exit_code: %d", containerJSON.State.Error, containerJSON.State.ExitCode))
					break
				}
				// Wait for 5 seconds before inspecting again
				<-time.After(5 * time.Second)
			}
		}
	}()

	// s.DockerClient.ContainerCreate(ctx, req.Container)
	return &container.Container{
		Name:        req.GetName(),
		Status:      "created",
		EnvVars:     req.GetEnvVars(),
		ContainerId: res.ID,
	}, nil
}

// GetContainer gets a container
func (s *service) GetContainer(ctx context.Context, req *container.GetContainerRequest) (*container.Container, error) {
	// inspect the container
	containerJSON, err := s.DockerClient.ContainerInspect(ctx, req.GetContainerId())
	if err != nil {
		log.Debug().Err(err).Msg("Error inspecting container")
		return nil, status.Error(codes.Internal, err.Error())
	}
	envVars := []*container.EnvVar{}
	for _, env := range containerJSON.Config.Env {
		split := strings.Split(env, "=")
		envVars = append(envVars, &container.EnvVar{
			Key:   split[0],
			Value: split[1],
		})
	}
	status := containerJSON.State.Status
	if containerJSON.State.Status == "running" || containerJSON.State.Status == "exited" && containerJSON.State.ExitCode == 0 {
		status = containerJSON.State.Status
	} else if containerJSON.State.Status == "exited" {
		status = fmt.Sprintf("failed: %s, exit_code: %d", containerJSON.State.Error, containerJSON.State.ExitCode)
	}
	return &container.Container{
		Name:        containerJSON.Name,
		Status:      status,
		EnvVars:     envVars,
		ContainerId: containerJSON.ID,
	}, nil
}

// DeleteContainer deletes a container
func (s *service) DeleteContainer(ctx context.Context, req *container.DeleteContainerRequest) (*container.DeleteContainerResponse, error) {
	// remove the container
	if err := s.DockerClient.ContainerRemove(ctx, req.GetContainerId(), dc.RemoveOptions{}); err != nil {
		log.Debug().Err(err).Msg("Error removing container")
		return nil, status.Error(codes.Internal, err.Error())
	}
	// remove the container from memory store
	memoryStoreKey := fmt.Sprintf("%s:%s", req.GetContainerId(), req.GetName())
	if err := s.MemoryStore.DeleteData(memoryStoreKey); err != nil {
		log.Debug().Err(err).Msg("Error removing container from memory store")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &container.DeleteContainerResponse{
		Message: fmt.Sprintf("Container %s removed", req.GetContainerId()),
	}, nil
}
