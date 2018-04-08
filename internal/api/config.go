package api

import "time"

type Config struct {
	TimeoutAfter time.Duration
	HackerNews   HackerNewsConfig   `toml:"HackerNews"`
	StoriesCache StoriesCacheConfig `toml:"StoriesCache"`
}

type HackerNewsConfig struct {
	BaseURL string
}

type StoriesCacheConfig struct {
	DefaultExpirationInMinutes time.Duration
	CleanupIntervalInMinutes   time.Duration
}
