package cache

import (
	"fmt"
	"github.com/donutloop/chn/internal/handler"
	"github.com/pkg/errors"
	"github.com/pmylund/go-cache"
	"time"
)

func NewStoriesCache(defaultExpirationInMinutes, cleanupIntervalInMinutes time.Duration) *StoriesCache {
	return &StoriesCache{
		cache: cache.New(defaultExpirationInMinutes*time.Minute, cleanupIntervalInMinutes*time.Minute),
	}
}

type StoriesCache struct {
	cache *cache.Cache
}

func (c *StoriesCache) GetStoriesBy(pageType string) ([]*handler.Story, error) {
	// check if we hit the cached stories for this page type
	value, found := c.cache.Get(pageType)
	if found {
		// check if valid stories
		s, ok := value.([]*handler.Story)
		if ok {
			return s, nil
		}
		return nil, errors.Errorf("element isn't a stories element (%#v)", s)
	}
	return nil, NewStoriesNotFoundError(pageType)
}

// SetStoriesBy save the stories into cache for this page type
func (c *StoriesCache) SetStoriesBy(pageType string, stories []*handler.Story) {
	c.cache.Set(pageType, stories, cache.DefaultExpiration)
}

func NewStoriesNotFoundError(pageType string) *StoriesNotFoundError {
	return &StoriesNotFoundError{pageType: pageType}
}

type StoriesNotFoundError struct {
	pageType string
}

func (err *StoriesNotFoundError) Error() string {
	return fmt.Sprintf("stories for %s not found", err.pageType)
}
