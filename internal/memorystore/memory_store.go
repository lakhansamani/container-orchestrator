package memorystore

import (
	"errors"
	"os"

	log "github.com/rs/zerolog/log"

	"github.com/lakhansamani/container-orchestrator/internal/memorystore/providers"
	"github.com/lakhansamani/container-orchestrator/internal/memorystore/providers/redis"
)

// NewMemoryStore initializes the memory store
func NewMemoryStore(redisURL string) (providers.MemoryStoreProvider, error) {
	// If redis url is not set throw an error
	if redisURL == "" {
		// Get the url from env
		redisURL = os.Getenv("REDIS_URL")
		if redisURL == "" {
			return nil, errors.New("redis url is not set")
		}
	}
	log.Info().Msg("Using redis store to save sessions")
	memoryStoreProvider, err := redis.NewRedisProvider(redisURL)
	if err != nil {
		return nil, err
	}
	return memoryStoreProvider, nil
}
