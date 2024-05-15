package memorystore

import (
	"errors"

	log "github.com/rs/zerolog/log"

	"github.com/lakhansamani/container-orchestrator/internal/memorystore/providers"
	"github.com/lakhansamani/container-orchestrator/internal/memorystore/providers/redis"
)

// NewMemoryStore initializes the memory store
func NewMemoryStore(redisURL string) (providers.MemoryStoreProvider, error) {
	// If redis url is not set throw an error
	if redisURL == "" {
		return nil, errors.New("redis url is not set")
	}
	log.Info().Msg("Using redis store to save sessions")
	memoryStoreProvider, err := redis.NewRedisProvider(redisURL)
	if err != nil {
		return nil, err
	}
	return memoryStoreProvider, nil
}
