package service

import (
	"github.com/donutloop/chn/internal/api"
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

func NewAPIService(port int, config *api.Config) *APIService {
	return &APIService{
		port:   port,
		config: config,
	}
}

type APIService struct {
	port   int
	config *api.Config
}

func (s *APIService) Init() error {

	// components ...
	hn := client.NewHackerNews(s.config.HackerNews.BaseURL, s.config.TimeoutAfter)
	github := client.NewGithub(s.config.Github.BaseURL, s.config.TimeoutAfter)

	storiesCache := cache.NewStoriesCache(s.config.StoriesCache.DefaultExpirationInMinutes, s.config.StoriesCache.CleanupIntervalInMinutes)

	// router and middleware ...
	r := chi.NewRouter()

	r.Use(
		middleware.DefaultLogger,
		middleware.Timeout(s.config.TimeoutAfter*time.Second),
		middleware.Recoverer,
	)

	// services ...
	storiesService := NewStoriesService(hn, storiesCache, github)

	// routes ...
	r.Method(http.MethodPost, "/xservice/service.chn.StoryService/Stories", handler.NewStoryServiceServer(storiesService, nil, log.Errorf))
	handler.FileServer(r, "/static", http.Dir("../../static"))
	r.Get("/", handler.File("index"))

	// start the server up on our port
	err := http.ListenAndServe(":"+strconv.Itoa(s.port), r)
	if err != nil {
		return err
	}

	return nil
}
