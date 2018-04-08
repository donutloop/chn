package service

import (
	"github.com/donutloop/chn/internal/cache"
	"github.com/donutloop/chn/internal/client"
	"github.com/donutloop/chn/internal/handler"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// baseURL is the URL for the hacker news API
var baseURL = "https://hacker-news.firebaseio.com/v0/"

func NewAPIService(port int) *APIService {
	return &APIService{port: port}
}

type APIService struct {
	port int
}

func (s *APIService) Init() error {

	hn := client.NewHackerNews(baseURL, 10)

	storiesCache := cache.NewStoriesCache(30, 10)

	http.Handle("/xservice/service.chn.StoryService/Stories", handler.NewStoryServiceServer(NewStoriesService(hn, storiesCache), nil, log.Errorf))

	// serve the favicon and logo files
	http.HandleFunc("/favicon.ico", handler.File("favicon"))
	http.HandleFunc("/logo.gif", handler.File("logo"))
	http.HandleFunc("/bundle.js", handler.File("bundle.js"))
	http.HandleFunc("/main.css", handler.File("main.css"))
	http.HandleFunc("/", handler.File("index"))

	// start the server up on our port
	err := http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
	if err != nil {
		return err
	}

	return nil
}
