// codehn is a hn clone that only displays posts from github
package main

// lots of imports means lots of time saved
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	// because the HN API is awkward and slow
	cache "github.com/pmylund/go-cache"
	"context"
	"github.com/gogo/protobuf/jsonpb"
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

// getStories if you couldn't guess it, gets the stories
func getStories(res *http.Response) ([]*Story, error) {

	// this is bad! we should limit the request body size
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// get all the story keys into a slice of ints
	var keys []int
	json.Unmarshal(body, &keys)

	// concurrency is cool, but needs to be limited
	semaphore := make(chan struct{}, 10)

	// how we know when all our goroutines are done
	wg := sync.WaitGroup{}

	// somewhere to store all the stories when we're done
	var stories []*Story

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
				return
			}
			defer res.Body.Close()

			s := &Story{}
			unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
			if err := unmarshaler.Unmarshal(res.Body, s); err != nil {
				log.Println(err)
				return
			}

			// parse the url
			u, err := url.Parse(s.Url)
			if err != nil {
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

// getStoriesFromType gets the different types of stories the API exposes
func getStoriesFromType(pageType string) ([]*Story, error) {
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
		return nil, errors.New("could not get " + pageType + " hacker news posts list")
	}

	defer res.Body.Close()
	s, err := getStories(res)
	if err != nil {
		return nil, errors.New("could not get " + pageType + " hacker news posts data")
	}

	return s, nil
}

type StoriesService struct {

}

// pageHandler returns a handler for the correct page type
func (service *StoriesService) Stories(ctx context.Context, req *StoryReq) (*StoryResp, error) {

	// we'll get all the stories
	var s []*Story

	// only because of shadowing
	var err error

	// know if we should use the cache
	var ok bool

	// check if we hit the cached stories for this page type
	value, found := cash.Get(req.Category)
	if found {
		// check if valid stories
		s, ok = value.([]*Story)
	}

	// if it's not or we didn't hit the cached stories
	if !ok {

		// get the stories from the API
		s, err = getStoriesFromType(req.Category)
		if err != nil {
			return nil, err
		}

		// set the cached stories for this page type
		cash.Set(req.Category, s, cache.DefaultExpiration)
	}

	resp := &StoryResp{
		Stories: s,
	}

	return resp, nil
}

// fileHandler serves a file like the favicon or logo
func fileHandler(file string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch file {
		case "favicon":
			http.ServeFile(w, r, "./favicon.ico")
		case "logo":
			http.ServeFile(w, r, "./logo.gif")
		case "index":
			http.ServeFile(w, r, "./static/html/index.html")
		case "bundle.js":
			http.ServeFile(w, r, "./static/js/bundle.js")
		case "main.css":
			http.ServeFile(w, r, "./static/css/main.css")
		default:
			w.WriteHeader(404)
			w.Write([]byte("file not found"))
		}
	}
}

func main() {

	// port 8080 is a good choice
	port := ":8080"

	http.Handle("/xservice/service.chn.StoryService/Stories", NewStoryServiceServer(new(StoriesService), nil))

	// serve the favicon and logo files
	http.HandleFunc("/favicon.ico", fileHandler("favicon"))
	http.HandleFunc("/logo.gif", fileHandler("logo"))
	http.HandleFunc("/bundle.js", fileHandler("bundle.js"))
	http.HandleFunc("/main.css", fileHandler("main.css"))
	http.HandleFunc("/", fileHandler("index"))

	// start the server up on our port
	log.Printf("Running on %s\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}
