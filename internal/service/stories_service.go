package service

import (
	"context"
	"github.com/donutloop/chn/internal/handler"
	cache "github.com/pmylund/go-cache"
	"net/url"
	"strings"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/donutloop/chn/internal/client"
)

// cash rules everything around me, get the money y'all
var cash *cache.Cache

func init() {
	// cash will have default expiration time of
	// 30 minutes and be swept every 10 minutes
	cash = cache.New(30*time.Minute, 10*time.Minute)
}

func NewStoriesService(hn *client.HackerNews) *StoriesService {
	return &StoriesService{
		hn: hn,
	}
}

type StoriesService struct{
	hn *client.HackerNews
}

// pageHandler returns a handler for the correct page type
func (service *StoriesService) Stories(ctx context.Context, req *handler.StoryReq) (*handler.StoryResp, error) {

	// we'll get all the stories
	 s := make([]*handler.Story, 0)

	// only because of shadowing
	var err error

	// know if we should use the cache
	var ok bool

	// check if we hit the cached stories for this page type
	value, found := cash.Get(req.Category)
	if found {
		// check if valid stories
		s, ok = value.([]*handler.Story)
	}

	// if it's not or we didn't hit the cached stories
	if !ok {

		// get the stories from the API
		s, err = service.getStoriesFromType(req.Category)
		if err != nil {
			log.WithError(err).Error("error get stories")
			return nil, err
		}

		// set the cached stories for this page type
		cash.Set(req.Category, s, cache.DefaultExpiration)
	}

	resp := &handler.StoryResp{
		Stories: s,
	}

	return resp, nil
}

// getStoriesFromType gets the different types of stories the API exposes
func (service *StoriesService) getStoriesFromType(pageType string) ([]*handler.Story, error) {
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

	stories, err := service.getStories(codes)
	if err != nil {
		return nil, err
	}

	return stories, nil
}

// getStories if you couldn't guess it, gets the stories
func (service *StoriesService) getStories(codes []int) ([]*handler.Story, error) {

	// concurrency is cool, but needs to be limited
	semaphore := make(chan struct{}, 10)

	// how we know when all our goroutines are done
	wg := sync.WaitGroup{}

	// somewhere to store all the stories when we're done
	stories := make([]*handler.Story, 0)

	// go over all the stories
	for _, code := range codes {

		// stop when we have 30 stories
		if len(stories) >= 30 {
			break
		}

		// sleep to avoid rate limiting from API
		time.Sleep(10 * time.Millisecond)

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
					Url: p.Url,
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
