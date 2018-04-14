package api

import "time"

type Config struct {
	ENV 		 string
	TimeoutAfter time.Duration
	HackerNews   HackerNewsConfig   `toml:"HackerNews"`
	StoriesCache StoriesCacheConfig `toml:"StoriesCache"`
	Github       GithubConfig       `toml:"Github"`
	Storage      StorageConfig      `toml:"Storage"`
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