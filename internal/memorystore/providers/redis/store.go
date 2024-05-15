package redis

import (
	"github.com/rs/zerolog/log"
)

// SetData sets the data in redis store.
func (c *memoryStoreProvider) SetData(key, value string) error {
	err := c.store.Set(c.ctx, key, value, -1).Err()
	if err != nil {
		log.Debug().Err(err).Msg("Error saving data to redis")
		return err
	}
	return nil
}

// GetData gets the data from redis store.
func (c *memoryStoreProvider) GetData(key string) (string, error) {
	data, err := c.store.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

// DeleteData deletes the data from redis store.
func (c *memoryStoreProvider) DeleteData(key string) error {
	err := c.store.Del(c.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
