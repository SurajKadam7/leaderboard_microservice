package config

import (
	"fmt"
	"time"

	"github.com/surajkadam/youtube_assignment/cache"
)

func New(key string, cache cache.Cache) *Config {
	return &Config{
		key:   key,
		Cache: cache,
	}
}

type Config struct {
	key   string
	Cache cache.Cache
}

func (c *Config) LifetimeKey() string {
	return c.key
}

func (c *Config) DayKey() string {
	y, m, d := time.Now().Date()
	key := fmt.Sprintf("%s:%d-%d-%d", c.key, y, m, d)
	return key
}
