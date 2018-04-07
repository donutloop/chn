package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/donutloop/chn/internal/handler"
	"github.com/golang/protobuf/jsonpb"
	"github.com/pkg/errors"
	cache "github.com/pmylund/go-cache"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
)

// baseURL is the URL for the hacker news API
var baseURL = "https://hacker-news.firebaseio.com/v0/"

// cash rules everything around me, get the money y'all
var cash *cache.Cache

func init() {
	// cash will have default expiration time of
	// 30 minutes and be swept every 10 minutes
	cash = cache.New(30*time.Minute, 10*time.Minute)
}

func NewStoriesService() *StoriesService {
	return new(StoriesService)
}

type StoriesService struct{}

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
	var url string
	switch pageType {
	case "best":
		url = baseURL + "beststories.json"
	case "new":
		url = baseURL + "newstories.json"
	case "show":
		url = baseURL + "showstories.json"
	default:
		url = baseURL + "topstories.json"
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Errorf("could not get %s hacker news posts list", pageType)
	}

	defer res.Body.Close()
	s, err := service.getStories(res)
	if err != nil {
		return nil, errors.Errorf("could not get %s hacker news posts data", pageType)
	}

	return s, nil
}

// getStories if you couldn't guess it, gets the stories
func (service *StoriesService) getStories(res *http.Response) ([]*handler.Story, error) {

	// this is bad! we should limit the request body size
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// get all the story keys into a slice of ints
	keys := make([]int, 0)
	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, err
	}

	// concurrency is cool, but needs to be limited
	semaphore := make(chan struct{}, 10)

	// how we know when all our goroutines are done
	wg := sync.WaitGroup{}

	// somewhere to store all the stories when we're done
	stories := make([]*handler.Story, 0)

	// go over all the stories
	for _, key := range keys {

		// stop when we have 30 stories
		if len(stories) >= 30 {
			break
		}

		// sleep to avoid rate limiting from API
		time.Sleep(10 * time.Millisecond)

		// in a goroutine with the story key
		go func(storyKey int) {

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

			// get the story with reckless abandon for errors
			keyURL := fmt.Sprintf(baseURL+"item/%d.json", storyKey)
			res, err := http.Get(keyURL)
			if err != nil {
				log.WithError(err).Error("error get stories")
				return
			}
			defer res.Body.Close()

			s := &handler.Story{}
			unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
			if err := unmarshaler.Unmarshal(res.Body, s); err != nil {
				log.WithError(err).Error("error get stories")
				return
			}

			// parse the url
			u, err := url.Parse(s.Url)
			if err != nil {
				log.WithError(err).Error("error get stories")
				return
			}

			// get the hostname from the url
			h := u.Hostname()

			// check if it's from github or gitlab before adding to stories
			if strings.Contains(h, "github") || strings.Contains(h, "gitlab") {
				s.DomainName = h
				stories = append(stories, s)
			}

		}(key)
	}

	// wait for all the goroutines
	wg.Wait()

	return stories, nil
}
