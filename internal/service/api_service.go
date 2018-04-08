package service

import (
	"github.com/donutloop/chn/internal/cache"
	"github.com/donutloop/chn/internal/client"
	"github.com/donutloop/chn/internal/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
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

	// components ...
	hn := client.NewHackerNews(baseURL, 10)

	storiesCache := cache.NewStoriesCache(30, 10)

	// router and middleware ...
	r := chi.NewRouter()

	r.Use(
		middleware.DefaultLogger,
		middleware.Timeout(15*time.Second),
		middleware.Recoverer,
	)

	// routes ...
	r.Method(http.MethodPost, "/xservice/service.chn.StoryService/Stories", handler.NewStoryServiceServer(NewStoriesService(hn, storiesCache), nil, log.Errorf))
	handler.FileServer(r, "/static", http.Dir("../../static"))
	r.Get("/", handler.File("index"))

	// start the server up on our port
	err := http.ListenAndServe(":"+strconv.Itoa(s.port), r)
	if err != nil {
		return err
	}

	return nil
}
