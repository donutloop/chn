package service

import (
	"context"
	"github.com/donutloop/chn/internal/cache"
	"github.com/donutloop/chn/internal/client"
	"github.com/donutloop/chn/internal/handler"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"sync"
	"github.com/donutloop/chn/internal/mediator"
	"github.com/donutloop/chn/internal/storage"
)

func NewStoriesService(hn *client.HackerNews, storiesCache *cache.StoriesCache, github *mediator.Github, st storage.Interface) *StoriesService {
	return &StoriesService{
		hn:           hn,
		storiesCache: storiesCache,
		gh:           github,
		st: st,
	}
}

type StoriesService struct {
	hn           *client.HackerNews
	gh           *mediator.Github
	storiesCache *cache.StoriesCache
	st           storage.Interface
}

// pageHandler returns a handler for the correct page type
func (service *StoriesService) Stories(ctx context.Context, req *handler.StoryReq) (*handler.StoryResp, error) {

	// we'll get all the stories
	stories := make([]*handler.Story, 0)
	var err error

	stories, err = service.storiesCache.GetStoriesBy(req.Category)
	if err != nil {
		if _, ok := err.(*cache.StoriesNotFoundError); ok {
			// get the stories from the API
			stories, err = service.getStoriesFromType(req.Category, req.Limit)
			if err != nil {
				log.WithError(err).Error("error get stories")
				return nil, err
			}
			service.storiesCache.SetStoriesBy(req.Category, stories)
		} else {
			return nil, err
		}
	}

	resp := &handler.StoryResp{
		Stories: stories,
	}

	return resp, nil
}

// getStoriesFromType gets the different types of stories the API exposes
func (service *StoriesService) getStoriesFromType(pageType string, limit int64) ([]*handler.Story, error) {
	var typ string
	switch pageType {
	case "best":
		typ = "beststories"
	case "new":
		typ = "newstories"
	case "show":
		typ = "showstories"
	default:
		typ = "topstories"
	}

	codes, err := service.hn.GetCodesForStory(typ)
	if err != nil {
		return nil, err
	}

	stories, err := service.getStories(codes, limit)
	if err != nil {
		return nil, err
	}

	return stories, nil
}

// getStories if you couldn't guess it, gets the stories
func (service *StoriesService) getStories(codes []int, limit int64) ([]*handler.Story, error) {

	// concurrency is cool, but needs to be limited
	semaphore := make(chan struct{}, 10)

	// how we know when all our goroutines are done
	wg := sync.WaitGroup{}

	// somewhere to store all the stories when we're done
	stories := make([]*handler.Story, 0)

	// go over all the stories
	for _, code := range codes {

		// stop when we have 30 stories
		if int64(len(stories)) >= limit {
			break
		}

		// in a goroutine with the story key
		go func(code int) {

			// add one to the wait group
			wg.Add(1)

			// add one to the semaphore
			semaphore <- struct{}{}

			// make sure this gets fired
			defer func() {

				// remove one from the wait group
				wg.Done()

				// remove one from the semaphore
				<-semaphore
			}()

			p, err := service.hn.GetPost(code)
			if err != nil {
				log.WithError(err).Error("error get stories")
				return
			}

			// parse the url
			u, err := url.Parse(p.Url)
			if err != nil {
				log.WithError(err).Error("error get stories")
				return
			}

			// get the hostname from the url
			h := u.Hostname()

			// check if it's from github or gitlab before adding to stories
			if strings.Contains(h, "github") || strings.Contains(h, "gitlab") {

				s := &handler.Story{
					Score: p.Score,
					Title: p.Title,
					Url:   p.Url,
				}

				if strings.Contains(h, "github") {
					var err error
					s.Langauges, err = service.gh.GetDataBy(p.Url)
					if err != nil {
						log.WithError(err).Error("error get stories")
					}
				}

				s.DomainName = h
				stories = append(stories, s)
			}

		}(code)
	}

	// wait for all the goroutines
	wg.Wait()

	return stories, nil
}
