package config

import "time"

type Config struct {
	TimeoutAfter time.Duration
	HackerNews   HackerNewsConfig   `toml:"HackerNews"`
	StoriesCache StoriesCacheConfig `toml:"StoriesCache"`
	Github       GithubConfig       `toml:"Github"`
	Storage      StorageConfig      `toml:"Storage"`
	StoriesStorage StoriesStorageConfig `toml:"StoriesStorage"`
}

type HackerNewsConfig struct {
	BaseURL string
}

type StoriesCacheConfig struct {
	DefaultExpirationInMinutes time.Duration
	CleanupIntervalInMinutes   time.Duration
}

type GithubConfig struct {
	BaseURL string
}

type StorageConfig struct {
	TimeoutAfter time.Duration
	Username string
	Password string
	Database string
	Handler  string
	Address  string
}

type StoriesStorageConfig struct {
	InitialInterval float64
	MaxInterval  float64
	Tries  			uint
}